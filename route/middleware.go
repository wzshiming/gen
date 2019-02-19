package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateMiddlewareFunction(midd *spec.Middleware) (err error) {
	name := g.getMiddlewareFunctionName(midd)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
// %s Is the middleware of %s
func %s(`, name, midd.Name, name)
	if midd.Type != nil {
		g.buf.WriteString("s *")
		g.Types(midd.Type)
		g.buf.WriteString(", ")
	}
	g.buf.WriteFormat(`w http.ResponseWriter, r *http.Request) (`)

	for i, resp := range midd.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		g.buf.WriteFormat("%s ", g.getVarName(resp.Name))
		g.Types(resp.Type)
	}
	g.buf.WriteString(`){
`)

	for _, req := range midd.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if req.Type == nil {
			continue
		}
		g.buf.WriteFormat("var %s ", g.getVarName(req.Name))
		g.Types(req.Type)
		g.buf.WriteString("\n")
	}

	err = g.generateCallExec(midd.Name, midd.PkgPath, midd.Type, midd.Requests, midd.Responses, true)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)

	return
}
