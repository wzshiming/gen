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

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
// %s Is the route of %s
func %s(`, name, oper.Name, name)
	if oper.Type != nil {
		g.buf.WriteString("s ")
		g.PtrTypes(oper.Type)
		g.buf.WriteString(", ")
	}
	g.buf.WriteFormat(`w http.ResponseWriter, r *http.Request) {
`)

	err = g.generateRequestsVar(oper.Requests, true)
	if err != nil {
		return err
	}

	err = g.generateResponsesVar(oper.Responses)
	if err != nil {
		return err
	}

	err = g.generateCallExec(oper.Name, oper.PkgPath, oper.Type, oper.Requests, oper.Responses, false)
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
				g.generateResponseBodyItem(resp)
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
