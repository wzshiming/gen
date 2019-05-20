package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateOperationFunction(oper *spec.Operation) (err error) {
	name := g.getOperationFunctionName(oper)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	err = g.generateFunctionDefine("route", name, oper.Name, oper.Type, nil, nil)
	if err != nil {
		return err
	}

	g.buf.WriteFormat(`{
`)
	err = g.generateRequestsVar(oper.Requests)
	if err != nil {
		return err
	}

	err = g.generateResponsesVar(oper.Responses)
	if err != nil {
		return err
	}

	errName, err := g.generateResponsesErrorName(oper.Responses)
	if err != nil {
		return err
	}

	err = g.generateCallExec(oper.Name, oper.Chain, oper.PkgPath, oper.Type, oper.Requests, oper.Responses, errName, false)
	if err != nil {
		return err
	}

	err = g.generateResponses(oper.Responses, "200", errName)
	if err != nil {
		return err
	}

	return
}
