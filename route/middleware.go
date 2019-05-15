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

	err = g.generateFunctionDefine("middleware", name, midd.Name, midd.Type)
	if err != nil {
		return err
	}

	g.buf.WriteFormat(`(`)
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
