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

	noCtx := true
	switch len(oper.Responses) {
	case 1:
		resp := oper.Responses[0]
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		typ := resp.Type
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		if typ.Kind != spec.Error {
			noCtx = false
		}
	case 0:
		// No action
	default:
		for _, resp := range oper.Responses {
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}
			if resp.In == "body" && resp.Content != "error" {
				g.generateResponseBodyItem(resp, errName)
				noCtx = false
				break
			}
		}
	}
	if noCtx {
		g.buf.WriteString(`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("null"))
`)
	}
	g.buf.WriteString(`
	return
}
`)
	return
}
