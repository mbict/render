package render

import (
	"net/http"

	"github.com/flosch/pongo2"
)

type pongo2Render struct {
	templateset *pongo2.TemplateSet
}

func NewPongo2Renderer(options ...Pongo2Options) Render {
	opt := preparePongo2Options(options)
	templateset := &pongo2.TemplateSet{}
	templateset.SetBaseDirectory(opt.Directory)
	templateset.Globals = opt.GlobalContext
	templateset.Debug = opt.DebugMode

	return &pongo2Render{
		templateset: templateset,
	}
}

func (r pongo2Render) Render(rw http.ResponseWriter, code int, data ...interface{}) error {

	name := data[0].(string)
	ctx, ok := data[1].(pongo2.Context)
	if !ok {
		ctx = pongo2.Context{}
	}
	ctx.Update(r.templateset.Globals)

	t := pongo2.Must(r.templateset.FromCache(name))

	rw.Header().Set("Content-Type", ContentHTML+"; charset=utf-8")
	rw.WriteHeader(code)
	t.ExecuteWriter(ctx, rw)
	return nil
}

type Pongo2Options struct {
	DebugMode     bool
	Directory     string
	GlobalContext pongo2.Context
}

func preparePongo2Options(options []Pongo2Options) Pongo2Options {
	var opt Pongo2Options
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}

	return opt
}
