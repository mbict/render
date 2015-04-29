package render

import (
	"encoding/json"
	"net/http"
)

const ContentJSON = "application/json"

type jsonRender struct{}

func (_ jsonRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {
	rw.Header().Set("Content-Type", ContentJSON+"; charset=utf-8")
	rw.WriteHeader(code)
	encoder := json.NewEncoder(rw)
	return encoder.Encode(data[0])
}
