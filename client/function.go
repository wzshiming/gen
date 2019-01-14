package client

import (
	"strings"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/namecase"
)

func (g *GenClient) GenerateFuncBody(oper *spec.Operation) (err error) {

	g.buf.WriteString("{\n")
	defer g.buf.WriteString(`
	return
}
`)
	g.GenerateRequests(oper)
	g.GenerateResponses(oper)
	return nil
}

func (g *GenClient) GenerateRequests(oper *spec.Operation) (err error) {

	g.buf.WriteString(`resp, err := Client.Clone()`)
	for _, v := range oper.Requests {
		req := v
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		switch req.In {
		case "security":
			// No action
		case "header":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`.
SetHead("%s", fmt.Sprint(%s))`, req.Name, g.GetVarName(req.Name))
		case "cookie":
			// TODO
		case "path":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`.
SetPath("%s", fmt.Sprint(%s))`, req.Name, g.GetVarName(req.Name))
		case "query":
			g.buf.WriteFormat(`.
SetQuery("%s", fmt.Sprint(%s))`, req.Name, g.GetVarName(req.Name))
		case "body":
			switch req.Content {
			case "json":
				g.buf.WriteFormat(`.
SetJSON(%s)`, g.GetVarName(req.Name))
			case "xml":
				g.buf.WriteFormat(`.
SetXML(%s)`, g.GetVarName(req.Name))
			case "file", "image":
				g.buf.WriteFormat(`.
SetBody(%s)`, g.GetVarName(req.Name))
			}
		}
	}
	g.buf.WriteFormat(`.
%s("%s")`, namecase.ToUpperHump(strings.SplitN(oper.Method, ",", 2)[0]), oper.Path)

	g.GenerateErrror(oper.Responses)
	return nil
}

func (g *GenClient) GenerateResponses(oper *spec.Operation) (err error) {

	g.buf.WriteString(`
	switch code := resp.StatusCode(); code {
`)

	for _, resp := range oper.Responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if resp.Code == "" {
			continue
		}
		g.buf.WriteFormat("case %s:", resp.Code)
		switch resp.Content {
		case "json":
			g.buf.AddImport("", "encoding/json")
			g.buf.WriteFormat(`
	err = json.Unmarshal(resp.Body(),&%s)
`, g.GetVarName(resp.Name))
		case "xml":
			g.buf.AddImport("", "encoding/xml")
			g.buf.WriteFormat(`
	err = xml.Unmarshal(resp.Body(),&%s)
`, g.GetVarName(resp.Name))
		case "error":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`
	%s = fmt.Errorf(string(resp.Body()))
`, g.GetVarName(resp.Name))
		}
		// TODO
	}
	g.buf.AddImport("", "net/http")
	g.buf.WriteString(`default:
		err = fmt.Errorf("Undefined code %d %s", code, http.StatusText(code))
	}
`)

	g.GenerateErrror(oper.Responses)
	return nil
}

func (g *GenClient) GenerateErrror(resps []*spec.Response) (err error) {
	g.buf.WriteString(`
if err != nil {
	return `)
	for i, resp := range resps {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if i != len(resps)-1 {
			g.TypesZero(resp.Type)
		} else {
			g.buf.WriteString("err")
		}
	}
	g.buf.WriteString(`
}
`)
	return
}
