package render

import (
	"fmt"
	"net/http"
)

const ContentPlain = "text/plain"

type plainRender struct{}

func (_ plainRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {
	rw.Header().Set("Content-Type", ContentPlain+"; charset=utf-8")
	rw.WriteHeader(code)
	format := data[0].(string)
	args := data[1].([]interface{})
	var err error
	if len(args) > 0 {
		_, err = rw.Write([]byte(fmt.Sprintf(format, args...)))
	} else {
		_, err = rw.Write([]byte(format))
	}
	return err
}
