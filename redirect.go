package render

import "net/http"

type redirectRender struct{}

func (_ redirectRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {
	rw.Header().Set("Location", data[0].(string))
	rw.WriteHeader(code)
	return nil
}
