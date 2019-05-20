package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateResponsesErrorName(resps []*spec.Response) (string, error) {
	for i := 0; i != len(resps); i++ {
		resp := resps[len(resps)-i-1]
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if resp.Type == nil {
			continue
		}
		if resp.Type.Kind == spec.Error {
			return g.getVarName(resp.Name, resp.Type), nil
		}
	}
	g.buf.WriteFormat("var err error")
	return "err", nil
}

func (g *GenRoute) generateResponsesVar(resps []*spec.Response) error {

	for _, resp := range resps {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if resp.Type == nil {
			continue
		}

		g.buf.WriteFormat("var %s ", g.getVarName(resp.Name, resp.Type))
		g.Types(resp.Type)
		g.buf.WriteString("\n")
	}

	return nil
}

func (g *GenRoute) generateResponse(resp *spec.Response, errName string) error {
	if resp.Ref != "" {
		resp = g.api.Responses[resp.Ref]
	}
	switch resp.In {
	case "header":
		return g.generateResponseHeader(resp)
	case "body":
		return g.generateResponseBody(resp, errName)
	}
	return nil
}

func (g *GenRoute) generateResponseHeader(resp *spec.Response) error {
	g.buf.AddImport("", "fmt")
	g.buf.WriteFormat(`
	w.Header().Set("%s",fmt.Sprint(%s))
`, resp.Name, g.getVarName(resp.Name, resp.Type))
	return nil
}

func (g *GenRoute) generateResponseBody(resp *spec.Response, errName string) error {
	text := ""
	if i, err := strconv.Atoi(resp.Code); err == nil {
		text = http.StatusText(i)
	}
	g.buf.WriteFormat(`
	// Response code %s %s for %s.
	if %s != `, resp.Code, text, resp.Name, g.getVarName(resp.Name, resp.Type))
	g.TypesZero(resp.Type)
	g.buf.WriteString(`{`)
	g.generateResponseBodyItem(resp.Name, resp.Type, resp.Code, resp.Content, errName, false)
	g.buf.WriteString(`return
}
`)
	return nil
}

func (g *GenRoute) generateResponseBodyItem(respName string, respType *spec.Type, respCode string, respContent string, errName string, noFmtErr bool) error {
	contentType := ""
	name := g.getVarName(respName, respType)

	switch respContent {
	case "json":
		g.buf.AddImport("", "encoding/json")
		contentType = "\"application/json; charset=utf-8\""
		g.buf.WriteFormat(`
	var _%s []byte
	_%s, %s = json.Marshal(%s)`, name, name, errName, name)
		g.generateResponseError(errName, "500", noFmtErr)
	case "xml":
		g.buf.AddImport("", "encoding/xml")
		contentType = "\"application/xml; charset=utf-8\""
		g.buf.WriteFormat(`
	var _%s []byte
	_%s, %s = xml.Marshal(%s)`, name, name, errName, name)
		g.generateResponseError(errName, "500", noFmtErr)
	case "error":
		g.generateResponseErrorReturn(name, respCode, noFmtErr)
		return nil
	default:
		typ := respType
		if typ.Attr.Has(spec.AttrReader) {
			g.buf.AddImport("", "io/ioutil")
			g.buf.WriteFormat(`
	var _%s []byte
	_%s, %s = ioutil.ReadAll(%s)
`, name, name, errName, name)
		} else if typ.Attr.Has(spec.AttrTextMarshaler) {
			g.buf.WriteFormat(`
	var _%s []byte
	_%s, %s = %s.MarshalText()
`, name, name, errName, name)
		} else if typ.Kind == spec.Slice && typ.Elem.Kind == spec.Byte {
			g.buf.WriteFormat(`
	var _%s []byte
	_%s = %s
`, name, name, name)
		} else {
			g.buf.AddImport("", "unsafe")
			g.buf.WriteFormat(`
	var _%s []byte
	var __%s string
	__%s = fmt.Sprint(%s)
	_%s = *(*[]byte)(unsafe.Pointer(&__%s))
`, name, name, name, name, name, name)
		}

		switch respContent {
		case "":
			contentType = "\"text/plain; charset=utf-8\""
		case "file":
			g.buf.AddImport("", "net/http")
			contentType = fmt.Sprintf("http.DetectContentType(_%s)", name)
		default:
			contentType = strconv.Quote(respContent)
		}

	}

	g.buf.WriteFormat(`
	w.Header().Set("Content-Type", %s)`, contentType)

	g.buf.WriteFormat(`
	w.WriteHeader(%s)`, respCode)

	switch respContent {
	case "xml":
		g.buf.AddImport("", "io")
		g.buf.WriteFormat(`
	io.WriteString(w, xml.Header)`)
	}

	g.buf.WriteFormat(`
	w.Write(_%s)
`, name)

	return nil
}
