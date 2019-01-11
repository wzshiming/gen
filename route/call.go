package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateCall(name, pkgpath string, typ *spec.Type, requests []*spec.Request, responses []*spec.Response, onErr bool) (err error) {

	for _, req := range requests {
		req := req
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if req.Content == "formdata" {
			g.buf.WriteFormat(`
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
`)
			break
		}
	}

	for _, req := range requests {
		err = g.GenerateRequest(req)
		if err != nil {
			return err
		}
	}

	if typ != nil {
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		g.buf.WriteFormat(`
	// Call %s %s.%s.
`, pkgpath, typ.Name, name)
	} else {
		g.buf.WriteFormat(`
	// Call %s %s.
`, pkgpath, name)
	}

	for i, resp := range responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		g.buf.WriteString(g.GetVarName(resp.Name))
	}
	if len(responses) != 0 {
		g.buf.WriteString(" = ")
	}
	if typ != nil {
		g.buf.WriteString("s.")
	} else {
		g.PkgPath(pkgpath)
	}

	g.buf.WriteFormat("%s(", name)
	for i, req := range requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if req.In == "none" {
			switch req.Ident {
			case "*net/http.Request":
				g.buf.WriteString("r")
			case "net/http.ResponseWriter":
				g.buf.WriteString("w")
			case "net/http.File":
				g.buf.WriteString(g.GetVarName(req.Name))
			}
		} else {
			g.buf.WriteString(g.GetVarName(req.Name))
		}
	}
	g.buf.WriteString(")\n")

	for _, resp := range responses {
		if onErr {
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}
			if resp.Name != "err" {
				continue
			}
		}
		err = g.GenerateResponse(resp)
		if err != nil {
			return err
		}
	}

	return nil
}
