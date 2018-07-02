package client

import (
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
)

func (g *GenClient) GenerateSecuritys() (err error) {
	for _, secu := range g.api.Securitys {
		err = g.GenerateSecurity(secu)
		if err != nil {
			return err
		}
		err = g.GenerateSecurityBody(secu)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GenClient) GenerateSecurity(oper *spec.Security) (err error) {
	g.buf.WriteString(utils.CommentLine(oper.Description))
	g.buf.WriteString("func ")

	if oper.Type != nil {
		g.buf.WriteByte('(')
		err = g.Types(oper.Type)
		if err != nil {
			return err
		}
		g.buf.WriteByte(')')
	}
	g.buf.WriteString(utils.GetName(oper.Name))
	g.buf.WriteByte('(')
	reqs := []*spec.Request{}
	for _, req := range oper.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		switch req.In {
		case "security":
			// No action
		case "header", "path", "query", "body":
			reqs = append(reqs, req)
		}
	}
	for i, req := range reqs {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		err = g.GenerateParameterRequests(req)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(")")
	return
}

func (g *GenClient) GenerateSecurityBody(oper *spec.Security) (err error) {
	g.buf.WriteString(`{
Client = Client`)
	defer g.buf.WriteString(`
}`)

	for _, req := range oper.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		switch req.In {
		case "security":
			// No action
		case "header":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`.
SetHeader("%s", fmt.Sprint(_%s))
`, req.Name, req.Name)
		case "cookie":
			// TODO
		case "path":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`.
SetPath("%s", fmt.Sprint(_%s))
`, req.Name, req.Name)
		case "query":
			g.buf.WriteFormat(`.
SetQuery("%s", fmt.Sprint(_%s))
`, req.Name, req.Name)
		case "body":
			// No action
		}
	}
	return nil
}
