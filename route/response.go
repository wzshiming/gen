package route

import (
	"net/http"
	"strconv"

	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateResponse(resp *spec.Response) error {
	if resp.Ref != "" {
		resp = g.api.Responses[resp.Ref]
	}
	switch resp.In {
	case "header":
		return g.GenerateResponseHeader(resp)
	case "body":
		return g.GenerateResponseBody(resp)
	}
	return nil
}

func (g *GenRoute) GenerateResponseHeader(resp *spec.Response) error {
	g.buf.AddImport("", "fmt")
	g.buf.WriteFormat(`
	w.Header().Set("%s",fmt.Sprint(%s))
`, resp.Name, g.GetVarName(resp.Name))
	return nil
}

func (g *GenRoute) GenerateResponseBody(resp *spec.Response) error {
	text := ""
	if i, err := strconv.Atoi(resp.Code); err == nil {
		text = http.StatusText(i)
	}
	g.buf.WriteFormat(`
	// Response code %s %s for %s.
	if %s != `, resp.Code, text, resp.Name, g.GetVarName(resp.Name))
	g.TypesZero(resp.Type)
	g.buf.WriteString(`{`)
	g.GenerateResponseBodyItem(resp)
	g.buf.WriteString(`return
}
`)
	return nil

}

func (g *GenRoute) GenerateResponseError() error {
	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
`)
	return nil
}

func (g *GenRoute) GenerateResponseBodyItem(resp *spec.Response) error {
	contentType := ""
	name := g.GetVarName(resp.Name)
	switch resp.Content {
	case "json":
		g.buf.AddImport("", "encoding/json")
		contentType = "application/json; charset=utf-8"
		g.buf.WriteFormat(`
	var _%s []byte
	_%s, err = json.Marshal(%s)`, name, name, name)
		g.GenerateResponseError()
	case "xml":
		g.buf.AddImport("", "encoding/xml")
		contentType = "application/xml; charset=utf-8"
		g.buf.WriteFormat(`
	var _%s []byte
	_%s, err = xml.Marshal(%s)`, name, name, name)
		g.GenerateResponseError()
	case "error":
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	http.Error(w, %s.Error(), %s)
`, name, resp.Code)
		return nil
	default:
		typ := resp.Type
		if typ.IsTextMarshaler {
			g.buf.WriteFormat(`
	var _%s []byte
	_%s, err = %s.MarshalText()
`, name, name, name)
		} else if typ.Kind == spec.Slice && typ.Elem.Kind == spec.Byte {
			g.buf.WriteFormat(`
	var _%s []byte
	_%s = %s
`, name, name, name)
		} else {
			g.buf.AddImport("", "unsafe")
			g.buf.WriteFormat(`
	var _%s []byte
	__%s = fmt.Sprint(_%s)
	%s = *(*[]byte)(unsafe.Pointer(&__%s))
`, name, name, name, name, name)
		}

		contentType = resp.Content
		if contentType == "" {
			contentType = "text/plain; charset=utf-8"
		}
	}

	g.buf.WriteFormat(`
	w.Header().Set("Content-Type","%s")
	w.WriteHeader(%s)
	w.Write(_%s)
`, contentType, resp.Code, name)
	return nil
}
