package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateOperationCall(oper *spec.Operation) error {
	g.buf.WriteFormat(`
	// Call %s.
`, oper.Name)
	for i, resp := range oper.Responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		g.buf.WriteString("_" + resp.Name)
	}
	g.buf.WriteString(":= ")
	if oper.Type != nil {
		g.buf.WriteString("s.")
	} else {
		g.PkgPath(oper.PkgPath)
	}

	g.buf.WriteFormat("%s(", oper.Name)
	for i, req := range oper.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if req.In == "none" {
			switch req.Name {
			case "*net/http.Request":
				g.buf.WriteString("r")
			case "net/http.ResponseWriter":
				g.buf.WriteString("w")
			}
		} else {
			g.buf.WriteString("_" + req.Name)
		}
	}
	g.buf.WriteString(")\n")
	return nil
}

func (g *GenRoute) GenerateOperationFunction(oper *spec.Operation) (err error) {
	name := g.GetOperationFunctionName(oper)

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
// %s Is the route of %s
func %s(`, name, oper.Name, name)
	if oper.Type != nil {
		g.buf.WriteString("s *")
		g.Types(oper.Type)
		g.buf.WriteString(", ")
	}
	g.buf.WriteFormat(`w http.ResponseWriter, r *http.Request) {
`)

	for _, req := range oper.Requests {
		err = g.GenerateOperationRequest(req)
		if err != nil {
			return err
		}
	}
	err = g.GenerateOperationCall(oper)
	if err != nil {
		return err
	}
	for _, resp := range oper.Responses {
		err = g.GenerateResponse(resp)
		if err != nil {
			return err
		}
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
				g.GenerateResponseBodyItem(resp)
				noCtx = false
				break
			}
		}
	}
	if noCtx {
		g.buf.WriteString(`
	w.WriteHeader(204)
	w.Write(nil)
`)
	}

	g.buf.WriteString(`
		return
}
`)
	return
}

func (g *GenRoute) GenerateOperationRequest(req *spec.Request) error {
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}

	switch req.In {
	case "none":
	// No action

	case "security":
		secus := []*spec.Security{}
		for _, secu := range g.api.Securitys {
			if len(secu.Responses) == 0 {
				continue
			}
			resp := secu.Responses[0]
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}

			if resp.Name != req.Name {
				continue
			}
			secus = append(secus, secu)
		}
		switch len(secus) {
		case 0:
			g.buf.WriteFormat(`
// Permission verification undefined.
var _%s `, req.Name)
			g.Types(req.Type)
			g.buf.WriteFormat(`
`)
		case 1:
			secu := secus[0]
			name := g.GetSecurityFunctionName(secu)
			g.buf.WriteFormat(`
// Permission verification call %s.
_%s, err := %s(r)
if err != nil {
	http.Error(w, err.Error(), 403)
	return
}
`, secu.Name, req.Name, name)
		default:
			secu := secus[0]
			name := g.GetSecurityFunctionName(secu)
			g.buf.WriteFormat(`
// Permission verification call %s.
_%s, err := %s(r)`, secu.Name, req.Name, name)
			for _, secu := range secus[1:] {
				g.buf.WriteFormat(`
if err != nil {
	// Permission verification call %s.
	_%s, err = %s(r)`, secu.Name, req.Name, name)
			}

			for range secus[1:] {
				g.buf.WriteString(`
}
`)
			}

			g.buf.WriteString(`
if err != nil {
	http.Error(w, err.Error(), 403)
	return
}
`)
		}
	default:
		g.buf.WriteFormat(`
// Parsing %s.
_%s, err := %s(r)`, req.Name, req.Name, g.GetRequestFunctionName(req))
		g.buf.WriteString(`
if err != nil {
	http.Error(w, err.Error(), 500)
	return
}
`)
	}

	return nil
}
