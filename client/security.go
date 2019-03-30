package client

import (
	"sort"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
)

func (g *GenClient) generateSecuritys() (err error) {
	secuKey := make([]string, 0, len(g.api.Securitys))
	for k := range g.api.Securitys {
		secuKey = append(secuKey, k)
	}
	sort.Strings(secuKey)
	for _, k := range secuKey {
		secu := g.api.Securitys[k]
		err = g.generateSecurity(secu)
		if err != nil {
			return err
		}
		err = g.generateSecurityBody(secu)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GenClient) generateSecurity(oper *spec.Security) (err error) {
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
	g.buf.WriteString(g.getSecurityName(oper))
	g.buf.WriteByte('(')

	switch oper.Schema {
	case "apiKey":
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
			err = g.generateParameterRequests(req, "")
			if err != nil {
				return err
			}
		}
	case "basic":
		g.buf.WriteString("username string, password string")
	case "bearer":
		g.buf.WriteString("token string")
	}

	g.buf.WriteString(")")
	return
}

func (g *GenClient) generateSecurityBody(oper *spec.Security) (err error) {
	g.buf.WriteString(`{
Client = Client`)

	switch oper.Schema {
	case "apiKey":
		for _, req := range oper.Requests {
			if req.Ref != "" {
				req = g.api.Requests[req.Ref]
			}
			switch req.In {
			case "none":

			case "security":
				// No action
			case "header":
				g.buf.AddImport("", "fmt")
				g.buf.WriteFormat(`.
SetHeader("%s", fmt.Sprint(%s))
`, req.Name, g.getVarName(req.Name, req.Type))
			case "cookie":
				// TODO
			case "path":
				g.buf.AddImport("", "fmt")
				g.buf.WriteFormat(`.
SetPath("%s", fmt.Sprint(%s))
`, req.Name, g.getVarName(req.Name, req.Type))
			case "query":
				g.buf.WriteFormat(`.
SetQuery("%s", fmt.Sprint(%s))
`, req.Name, g.getVarName(req.Name, req.Type))
			case "body":
				// No action
			}
		}
	case "basic":
		g.buf.WriteString(`.
SetBasicAuth(username, password)
`)
	case "bearer":
		g.buf.WriteString(`.
SetAuthToken(token)
`)
	}
	g.buf.WriteString(`
	}
`)
	return nil
}
