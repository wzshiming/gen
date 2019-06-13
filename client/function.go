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
		case "header", "cookie", "path":
			name := g.getVarName(req.Name, req.Type)
			g.buf.WriteFormat(`var _%s string
`, name)
		case "query":
			name := g.getVarName(req.Name, req.Type)
			if g.explode && (req.Type.Kind == spec.Array || req.Type.Kind == spec.Slice) {
				g.buf.WriteFormat(`var _%s []string
`, name)
			} else {
				g.buf.WriteFormat(`var _%s string
`, name)
			}
		}
	}

	for _, v := range oper.Requests {
		req := v
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		switch req.In {
		case "header", "cookie", "path":
			name := g.getVarName(req.Name, req.Type)
			err = g.GenModel.ConvertFrom(name, "_"+name, req.Type)
			if err != nil {
				return err
			}
		case "query":
			name := g.getVarName(req.Name, req.Type)
			if req.Type.Kind == spec.Array || req.Type.Kind == spec.Slice {
				err = g.GenModel.ConvertFromMulti(name, "_"+name, req.Type, g.explode)
				if err != nil {
					return err
				}
			} else {
				err = g.GenModel.ConvertFrom(name, "_"+name, req.Type)
				if err != nil {
					return err
				}
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
			if g.explode && (req.Type.Kind == spec.Array || req.Type.Kind == spec.Slice) {
				g.buf.WriteFormat(`.
AddQuerys("%s", _%s)`, req.Name, g.getVarName(req.Name, req.Type))
			} else {
				g.buf.WriteFormat(`.
SetQuery("%s", _%s)`, req.Name, g.getVarName(req.Name, req.Type))
			}
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

func (g *GenClient) generateResponsesUnmarshal(resp *spec.Response, define bool) (err error) {
	if resp.Ref != "" {
		resp = g.api.Responses[resp.Ref]
	}
	if resp.Code == "" {
		return nil
	}
	name := g.getVarName(resp.Name, resp.Type)

	if define {
		g.buf.WriteFormat(`
	var %s `, name)
		err = g.Types(resp.Type)
		if err != nil {
			return err
		}
	} else {
		g.buf.WriteFormat("case %s:", resp.Code)
	}
	switch resp.Content {
	case "json":
		g.buf.AddImport("", "encoding/json")
		g.buf.WriteFormat(`
	err = json.Unmarshal(resp.Body(), &%s)
`, name)
		if define {
			g.buf.WriteFormat(`
	if err == nil {
		err = fmt.Errorf("%%+v", %s)
	}
`, name)
		}
	case "xml":
		g.buf.AddImport("", "encoding/xml")
		g.buf.WriteFormat(`
	err = xml.Unmarshal(resp.Body(), &%s)
`, name)
		if define {
			g.buf.WriteFormat(`
	if err == nil {
		err = fmt.Errorf("%%+v", %s)
	}
`, name)
		}
	case "error":
		g.buf.AddImport("", "fmt")
		for _, wrap := range g.api.Wrappings {
			if len(wrap.Responses) == 0 {
				continue
			}
			res := wrap.Responses[0]
			if res.Ref != "" {
				res = g.api.Responses[res.Ref]
			}
			if res.Content != resp.Content {
				return g.generateResponsesUnmarshal(res, true)
			}
		}

		g.buf.WriteFormat(`
	%s = fmt.Errorf(string(resp.Body()))
`, name)
	}

	return nil
}

func (g *GenClient) generateResponses(oper *spec.Operation) (err error) {
	g.buf.WriteString(`
	switch code := resp.StatusCode(); code {
`)

	for _, resp := range oper.Responses {
		err := g.generateResponsesUnmarshal(resp, false)
		if err != nil {
			return err
		}
	}
	g.buf.AddImport("", "net/http")
	g.buf.WriteString(`
	default:
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
