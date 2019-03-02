package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateCallExec(name, pkgpath string, typ *spec.Type, requests []*spec.Request, responses []*spec.Response, onErr bool) (err error) {

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

	err = g.generateCall(name, pkgpath, typ, requests, responses)
	if err != nil {
		return err
	}

	for _, resp := range responses {
		if onErr {
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}
			if resp.Name != "err" {
				continue
			}
		}
		err = g.generateResponse(resp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GenRoute) generateCall(name, pkgpath string, typ *spec.Type, requests []*spec.Request, responses []*spec.Response) (err error) {

	for _, req := range requests {
		err = g.generateRequest(req)
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

	if len(responses) != 0 {
		for i, resp := range responses {
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}
			if i != 0 {
				g.buf.WriteByte(',')
			}
			g.buf.WriteString(g.getVarName(resp.Name, resp.Type))
		}
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
		g.buf.WriteString(g.getVarName(req.Name, req.Type))
	}
	g.buf.WriteString(")\n")

	return nil
}
