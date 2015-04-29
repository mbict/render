package render

import (
	"encoding/json"
	"net/http"
)

const ContentJSONP = "application/javascript"

type jsonpRender struct {
	indent   bool
	callback string //default callback
}

func NewJsonpRenderer(callback string, indent bool) Render {
	return &jsonpRender{
		indent:   indent,
		callback: callback,
	}
}

func (j jsonpRender) Render(rw http.ResponseWriter, code int, data ...interface{}) error {
	rw.Header().Set("Content-Type", ContentJSONP+"; charset=utf-8")
	rw.WriteHeader(code)

	var (
		result   []byte
		err      error
		callback string
		v        interface{}
	)

	switch len(callback) {
	case 0:
		callback = j.callback
		v = nil

	case 1:
		callback = j.callback
		v = data[0]

	default:
		callback = data[0].(string)
		v = data[1]
	}

	if j.indent {
		result, err = json.MarshalIndent(v, "", "  ")
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}

	rw.Write([]byte(callback + "("))
	rw.Write(result)
	rw.Write([]byte(");"))

	// If indenting, append a new line.
	if j.indent {
		rw.Write([]byte("\n"))
	}
	return nil
}
