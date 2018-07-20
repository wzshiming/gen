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
w.Header().Set("%s",fmt.Sprint(_%s))
`, resp.Name, resp.Name)
	return nil
}

func (g *GenRoute) GenerateResponseBody(resp *spec.Response) error {
	text := ""
	if i, err := strconv.Atoi(resp.Code); err == nil {
		text = http.StatusText(i)
	}
	g.buf.WriteFormat(`
	// Response code %s %s for %s.
	if _%s != `, resp.Code, text, resp.Name, resp.Name)
	g.TypesZero(resp.Type)
	g.buf.WriteString(`{`)
	g.GenerateResponseBodyItem(resp)
	g.buf.WriteString(`
	return
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

	switch resp.Content {
	case "json":
		g.buf.AddImport("", "encoding/json")
		contentType = "application/json; charset=utf-8"
		g.buf.WriteFormat(`
	data, err := json.Marshal(_%s)`, resp.Name)
		g.GenerateResponseError()
	case "xml":
		g.buf.AddImport("", "encoding/xml")
		contentType = "application/xml; charset=utf-8"
		g.buf.WriteFormat(`
	data, err := xml.Marshal(_%s)`, resp.Name)
		g.GenerateResponseError()
	case "error":
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	http.Error(w, _%s.Error(), %s)
`, resp.Name, resp.Code)
		return nil
	default:
		contentType = "text/plain; charset=utf-8"
	}

	g.buf.WriteFormat(`
	w.Header().Set("Content-Type","%s")
	w.WriteHeader(%s)
	w.Write(data)
`, contentType, resp.Code)
	return nil
}
