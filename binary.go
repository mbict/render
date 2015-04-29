package render

import "net/http"

const ContentBinary = "application/octet-stream"

type binaryRender struct{}

func (_ binaryRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {
	rw.Header().Set("Content-Type", ContentBinary+"; charset=utf-8")
	rw.WriteHeader(code)
	_, err := rw.Write(data[0].([]byte))
	return err
}
