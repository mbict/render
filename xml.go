package render

import (
	"encoding/xml"
	"net/http"
)

const ContentXML = "text/xml"

type xmlRender struct{}

func (_ xmlRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {
	rw.Header().Set("Content-Type", ContentXML+"; charset=utf-8")
	rw.WriteHeader(code)
	encoder := xml.NewEncoder(rw)
	return encoder.Encode(data[0])
}
