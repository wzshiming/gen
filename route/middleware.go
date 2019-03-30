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
		g.buf.WriteString("s ")
		g.PtrTypes(midd.Type)
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
		g.buf.WriteFormat("%s ", g.getVarName(resp.Name, resp.Type))
		g.Types(resp.Type)
	}
	g.buf.WriteString(`) {
`)
	err = g.generateRequestsVar(midd.Requests, false)
	if err != nil {
		return err
	}

	err = g.generateCallExec(midd.Name, nil, midd.PkgPath, midd.Type, midd.Requests, midd.Responses, true)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)

	return
}
