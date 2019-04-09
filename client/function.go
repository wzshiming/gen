package client

import (
	"strings"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/namecase"
)

func (g *GenClient) generateFuncBody(oper *spec.Operation) (err error) {

	g.buf.WriteString(`{
`)
	err = g.generateRequests(oper)
	if err != nil {
		return err
	}
	err = g.generateResponses(oper)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)
	return nil
}

func (g *GenClient) generateRequests(oper *spec.Operation) (err error) {
	for _, v := range oper.Requests {
		req := v
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		switch req.In {
		case "header", "cookie", "path", "query":
			name := g.getVarName(req.Name, req.Type)
			g.buf.WriteFormat(`var _%s string
`, name)
		}
	}

	for _, v := range oper.Requests {
		req := v
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		switch req.In {
		case "header", "cookie", "path", "query":
			name := g.getVarName(req.Name, req.Type)
			err = g.GenModel.ConvertFrom(name, "_"+name, req.Type)
			if err != nil {
				return err
			}
		}
	}

	g.buf.WriteString(`
	resp, err := Client.Clone()`)
	for _, v := range oper.Requests {
		req := v
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		switch req.In {
		case "security":
			// No action
		case "header":
			g.buf.WriteFormat(`.
SetHeader("%s", _%s)`, req.Name, g.getVarName(req.Name, req.Type))
		case "cookie":
			// TODO
		case "path":
			g.buf.WriteFormat(`.
SetPath("%s", _%s)`, req.Name, g.getVarName(req.Name, req.Type))
		case "query":
			g.buf.WriteFormat(`.
SetQuery("%s", _%s)`, req.Name, g.getVarName(req.Name, req.Type))
		case "body":
			switch req.Content {
			case "json":
				g.buf.WriteFormat(`.
SetJSON(%s)`, g.getVarName(req.Name, req.Type))
			case "xml":
				g.buf.WriteFormat(`.
SetXML(%s)`, g.getVarName(req.Name, req.Type))
			case "file", "image":
				g.buf.WriteFormat(`.
SetBody(%s)`, g.getVarName(req.Name, req.Type))
			}
		}
	}
	g.buf.WriteFormat(`.
%s("%s")`, namecase.ToUpperHump(strings.SplitN(oper.Method, ",", 2)[0]), strings.Trim(oper.Path, "/"))

	g.generateErrror(oper.Responses)
	return nil
}

func (g *GenClient) generateResponses(oper *spec.Operation) (err error) {

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
`, g.getVarName(resp.Name, resp.Type))
		case "xml":
			g.buf.AddImport("", "encoding/xml")
			g.buf.WriteFormat(`
	err = xml.Unmarshal(resp.Body(),&%s)
`, g.getVarName(resp.Name, resp.Type))
		case "error":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`
	%s = fmt.Errorf(string(resp.Body()))
`, g.getVarName(resp.Name, resp.Type))
		}
		// TODO
	}
	g.buf.AddImport("", "net/http")
	g.buf.WriteString(`default:
		if code >= 400 {
			err = fmt.Errorf("Undefined code %d %s", code, http.StatusText(code))
		}
	}
`)

	g.generateErrror(oper.Responses)
	return nil
}

func (g *GenClient) generateErrror(resps []*spec.Response) (err error) {
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
