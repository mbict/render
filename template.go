package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Included helper functions for use when rendering html
var helperFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called with no layout defined")
	},
}

// Delims represents a set of Left and Right delimiters for HTML template rendering
type TemplateDelims struct {
	// Left delimiter, defaults to {{
	Left string
	// Right delimiter, defaults to }}
	Right string
}

type TemplateOptions struct {
	DebugMode  bool
	Directory  string
	Layout     string
	Extensions []string
	Funcs      []template.FuncMap
	Delims     TemplateDelims
}

type templateRender struct {
	options     TemplateOptions
	htmlOptions HTMLOptions
	template    *template.Template
}

type HTMLOptions struct {
	// Layout template name. Overrides Options.Layout.
	Layout string
}

func NewTemplateRenderer(options ...TemplateOptions) Render {
	opt := prepareTemplateOptions(options)
	return &templateRender{
		options:  opt,
		template: compileTemplate(opt),
	}
}

func (r *templateRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {

	var t *template.Template

	//recompile every request if in debug mode
	if r.options.DebugMode == true {
		t = compileTemplate(r.options)
	} else {
		t = template.Must(r.template.Clone())
	}

	//get options
	name := data[0].(string)
	binding := data[1]
	var opt HTMLOptions
	if len(data) > 2 {
		opt = data[2].(HTMLOptions)
	} else {
		opt = r.htmlOptions
	}

	if len(opt.Layout) > 0 {
		r.addYield(t, name, binding)
		name = opt.Layout
	}

	out, err := r.execute(t, name, binding)
	if err != nil {
		panic(err)
	}

	// template rendered fine, write out the result
	rw.Header().Set("Content-Type", ContentHTML+"; charset=utf-8")
	rw.WriteHeader(code)
	io.Copy(rw, out)
	return nil
}

func prepareTemplateOptions(options []TemplateOptions) TemplateOptions {
	var opt TemplateOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.Extensions) == 0 {
		opt.Extensions = []string{".tmpl"}
	}

	return opt
}

func compileTemplate(options TemplateOptions) *template.Template {
	dir := options.Directory
	t := template.New(dir)
	t.Delims(options.Delims.Left, options.Delims.Right)
	// parse an initial template in case we don't have any
	template.Must(t.Parse(""))

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(r)
		for _, extension := range options.Extensions {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := (r[0 : len(r)-len(ext)])
				tmpl := t.New(filepath.ToSlash(name))

				// add our funcmaps
				for _, funcs := range options.Funcs {
					tmpl.Funcs(funcs)
				}

				// Bomb out if parse fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				break
			}
		}

		return nil
	})

	return t
}

func (r *templateRender) execute(t *template.Template, name string, binding interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	return buf, t.ExecuteTemplate(buf, name, binding)
}

func (r *templateRender) addYield(t *template.Template, name string, binding interface{}) {
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf, err := r.execute(t, name, binding)
			// return safe html here since we are rendering our own template
			return template.HTML(buf.String()), err
		},
	}
	t.Funcs(funcs)
}

func (r *templateRender) prepareHTMLOptions(htmlOpt []HTMLOptions) HTMLOptions {
	if len(htmlOpt) > 0 {
		return htmlOpt[0]
	}

	return HTMLOptions{
		Layout: r.options.Layout,
	}
}
