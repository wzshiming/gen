package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateSecurityCall(secu *spec.Security) error {
	g.buf.WriteFormat(`
		// Call %s.
	`, secu.Name)
	g.buf.WriteString("return ")
	if secu.Type != nil {
		typ := secu.Type
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}

		g.buf.WriteFormat("%s.", GetGlobalVarName(typ.Name))
	}
	g.buf.WriteFormat("%s(", secu.Name)
	for i, req := range secu.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		g.buf.WriteString("_" + req.Name)
	}
	g.buf.WriteString(")\n")
	return nil
}

func (g *GenRoute) GenerateSecurityFunction(secu *spec.Security) (err error) {
	name := GetSecurityFunctionName(secu.Name)
	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
	// %s Is the security of %s
	func`, name, secu.Name)
	//	if secu.Type != nil {
	//		g.buf.WriteString("(s ")
	//		g.Types(secu.Type)
	//		g.buf.WriteString(")")
	//	}
	g.buf.WriteFormat(` %s(r *http.Request)`, name)

	g.buf.WriteString("(")
	for i, resp := range secu.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		g.buf.WriteFormat("%s ", resp.Name)
		g.Types(resp.Type)
	}
	g.buf.WriteString(")")

	g.buf.WriteString(`{
	`)
	defer g.buf.WriteString(`
	}
	`)

	for _, req := range secu.Requests {
		err = g.GenerateSecurityRequest(req)
		if err != nil {
			return err
		}
	}
	err = g.GenerateSecurityCall(secu)
	if err != nil {
		return err
	}
	return
}

func (g *GenRoute) GenerateSecurityRequest(req *spec.Request) error {
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}

	g.buf.WriteFormat(`
// Parsing %s.
_%s, err := %s(r)`, req.Name, req.Name, GetRequestFunctionName(req.Name, req.In))
	g.buf.WriteString(`
if err != nil {
	return `)

	g.buf.WriteString(`
}
`)
	return nil
}
