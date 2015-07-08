package render

import (
	"net/http"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type markdownRender struct{}

func (_ markdownRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {
	rw.Header().Set("Content-Type", ContentHTML+"; charset=utf-8")
	rw.WriteHeader(code)
	input := data[0].([]byte)

	unsafe := blackfriday.MarkdownCommon(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	_, err := rw.Write(html)
	return err
}
