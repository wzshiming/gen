package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateSecurityCall(secu *spec.Security) error {
	g.buf.WriteFormat(`
		// Call %s.
`, secu.Name)
	errFlag := false
	for i, resp := range secu.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if !errFlag && resp.Name == "err" {
			errFlag = true
		}
		g.buf.WriteFormat("%s ", resp.Name)
	}
	g.buf.WriteString(` = `)
	g.PkgPath(secu.PkgPath)
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
	g.buf.WriteString(`)
`)
	if errFlag {
		g.buf.WriteString(`
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
`)
	}
	g.buf.WriteString(`
		return
`)
	return nil
}

func (g *GenRoute) GenerateSecurityFunction(secu *spec.Security) (err error) {
	name := g.GetSecurityFunctionName(secu)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
	// %s Is the security of %s
	func %s(w http.ResponseWriter, r *http.Request) (`, name, secu.Name, name)

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
	g.buf.WriteString(`){
`)
	for _, req := range secu.Requests {
		err = g.GenerateOperationRequest(req)
		if err != nil {
			return err
		}
	}
	err = g.GenerateSecurityCall(secu)
	if err != nil {
		return err
	}

	g.buf.WriteString(`
}
`)

	return
}
