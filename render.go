package render

import "net/http"

//default renders
var (
	Markdown  = &markdownRender{}
	HtmlPlain = &htmlPlainRender{}
	Binary    = &binaryRender{}
	XML       = &xmlRender{}
	Plain     = &plainRender{}
	Redirect  = &redirectRender{}
	JSON      = &jsonRender{}
	JSONP     = &jsonpRender{
		indent:   true,
		callback: "callback",
	}
)

type Render interface {
	Render(http.ResponseWriter, int, ...interface{}) error
}
