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

	err = g.generateFunctionDefine("middleware", name, midd.Name, midd.Type, nil, midd.Responses)
	if err != nil {
		return err
	}

	g.buf.WriteString(`{
`)
	err = g.generateRequestsVar(midd.Requests)
	if err != nil {
		return err
	}

	errName, err := g.generateResponsesErrorName(midd.Responses)
	if err != nil {
		return err
	}

	err = g.generateCallExec(midd.Name, nil, midd.PkgPath, midd.Type, midd.Requests, midd.Responses, errName, true)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)

	return
}
