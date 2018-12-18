package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateMiddlewareCall(midd *spec.Middleware) error {
	g.buf.WriteFormat(`
		// Call %s.
`, midd.Name)
	errFlag := false
	for i, resp := range midd.Responses {
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
	g.PkgPath(midd.PkgPath)
	g.buf.WriteFormat("%s(", midd.Name)
	for i, req := range midd.Requests {
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
		http.Error(w, err.Error(), 400)
		return
	}
`)
	}
	g.buf.WriteString(`
	return
`)
	return nil
}

func (g *GenRoute) GenerateMiddlewareFunction(midd *spec.Middleware) (err error) {
	name := g.GetMiddlewareFunctionName(midd)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
	// %s Is the middleware of %s
	func %s(w http.ResponseWriter, r *http.Request) (`, name, midd.Name, name)

	for i, resp := range midd.Responses {
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
	for _, req := range midd.Requests {
		err = g.GenerateOperationRequest(req)
		if err != nil {
			return err
		}
	}
	err = g.GenerateMiddlewareCall(midd)
	if err != nil {
		return err
	}

	g.buf.WriteString(`
}
`)

	return
}
