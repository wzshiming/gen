package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateSecurityFunction(secu *spec.Security) (err error) {
	name := g.getSecurityFunctionName(secu)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
// %s Is the security of %s
func %s(`, name, secu.Name, name)
	if secu.Type != nil {
		g.buf.WriteString("s *")
		g.Types(secu.Type)
		g.buf.WriteString(", ")
	}
	g.buf.WriteFormat(`w http.ResponseWriter, r *http.Request) (`)

	for i, resp := range secu.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		g.buf.WriteFormat("%s ", g.getVarName(resp.Name, resp.Type))
		g.Types(resp.Type)
	}
	g.buf.WriteString(`){
`)

	for _, req := range secu.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if req.Type == nil {
			continue
		}
		g.buf.WriteFormat("var %s ", g.getVarName(req.Name, req.Type))
		g.Types(req.Type)
		g.buf.WriteString("\n")
	}

	err = g.generateCallExec(secu.Name, secu.PkgPath, secu.Type, secu.Requests, secu.Responses, true)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)

	return
}
