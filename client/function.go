package client

import (
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
SetHead("%s", fmt.Sprint(_%s))`, req.Name, req.Name)
		case "cookie":
			// TODO
		case "path":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`.
SetPath("%s", fmt.Sprint(_%s))`, req.Name, req.Name)
		case "query":
			g.buf.WriteFormat(`.
SetQuery("%s", fmt.Sprint(_%s))`, req.Name, req.Name)
		case "body":
			switch req.Content {
			case "json":
				g.buf.WriteFormat(`.
SetJSON(_%s)`, req.Name)
			case "xml":
				g.buf.WriteFormat(`.
SetXML(_%s)`, req.Name)
			}
		}
	}
	g.buf.WriteFormat(`.
%s("%s")`, namecase.ToUpperHump(oper.Method), oper.Path)

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
	err = json.Unmarshal(resp.Body(),&_%s)
`, resp.Name)
		case "xml":
			g.buf.AddImport("", "encoding/xml")
			g.buf.WriteFormat(`
	err = xml.Unmarshal(resp.Body(),&_%s)
`, resp.Name)
		case "error":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`
	_%s = fmt.Errorf(string(resp.Body()))
`, resp.Name)
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
