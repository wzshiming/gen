package route

import (
	"fmt"

	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateRequest(req *spec.Request) error {
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}

	vname := g.GetVarName(req.Name)
	switch req.In {
	case "none":
		// No action

	case "middleware":
		midds := []*spec.Middleware{}
		for _, midd := range g.api.Middlewares {
			if len(midd.Responses) == 0 {
				continue
			}
			resp := midd.Responses[0]
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}

			if resp.Name != req.Name {
				continue
			}
			midds = append(midds, midd)
		}
		switch len(midds) {
		case 0:
			g.buf.WriteFormat(`
// Permission middleware undefined %s.
`, req.Name)
		case 1:
			secu := midds[0]
			name := g.GetMiddlewareFunctionName(secu)
			g.buf.WriteFormat(`
// Permission middlewares call %s.
%s, err = %s(`, secu.Name, vname, name)
			if secu.Type != nil {
				g.buf.WriteString(`s, `)
			}
			g.buf.WriteString(`w, r)
if err != nil {
	return
}
`)
		}
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
var %s `, vname)
			g.Types(req.Type)
			g.buf.WriteFormat(`
`)
		case 1:
			secu := secus[0]
			name := g.GetSecurityFunctionName(secu)
			g.buf.WriteFormat(`
// Permission verification call %s.
%s, err = %s(`, secu.Name, vname, name)
			if secu.Type != nil {
				g.buf.WriteString(`s, `)
			}
			g.buf.WriteString(`w, r)
if err != nil {
	return
}
`)
		default:
			secu := secus[0]
			name := g.GetSecurityFunctionName(secu)
			g.buf.WriteFormat(`
// Permission verification call %s.
%s, err = %s(`, secu.Name, vname, name)
			if secu.Type != nil {
				g.buf.WriteString(`s, `)
			}
			g.buf.WriteString(`w, r)`)
			for _, secu := range secus[1:] {
				g.buf.WriteFormat(`
if err != nil {
	// Permission verification call %s.
	%s, err = %s(`, secu.Name, vname, name)
				if secu.Type != nil {
					g.buf.WriteString(`s, `)
				}
				g.buf.WriteString(`w, r)`)
			}

			for range secus[1:] {
				g.buf.WriteString(`
}
`)
			}

			g.buf.WriteString(`
if err != nil {
	return
}
`)
		}
	default:
		g.buf.WriteFormat(`
// Parsing %s.
%s, err = %s(w, r)`, req.Name, vname, g.GetRequestFunctionName(req))
		g.buf.WriteString(`
if err != nil {
	return
}
`)
	}

	return nil
}

func (g *GenRoute) GenerateRequestFunction(req *spec.Request) error {
	g.buf.AddImport("", "net/http")

	name := g.GetRequestFunctionName(req)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	g.buf.WriteFormat(`
// %s Parsing the %s for of %s
func %s(w http.ResponseWriter, r *http.Request) (%s `, name, req.In, req.Name, name, g.GetVarName(req.Name))
	g.Types(req.Type)
	g.buf.WriteString(`, err error) {
`)
	err := g.GenerateRequestVar(req)
	if err != nil {
		return err
	}
	g.buf.WriteString(`

	return
}`)
	return nil
}

func (g *GenRoute) GenerateRequestVar(req *spec.Request) error {

	name := g.GetVarName(req.Name)
	switch req.In {
	case "body":

		switch req.Content {
		case "json":
			g.buf.AddImport("", "io/ioutil")
			g.buf.AddImport("", "encoding/json")
			g.buf.WriteFormat(`
	var _%s []byte
	_%s, err = ioutil.ReadAll(r.Body)
	if err == nil {
		r.Body.Close()
		err = json.Unmarshal(_%s, &%s)
	}
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`, name, name, name, name)
		case "xml":
			g.buf.AddImport("", "io/ioutil")
			g.buf.AddImport("", "encoding/xml")
			g.buf.WriteFormat(`
	var _%s []byte
	_%s, err = ioutil.ReadAll(r.Body)
	if err == nil {
		r.Body.Close()
		err = xml.Unmarshal(_%s, &%s)
	}
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`, name, name, name, name)
		case "formdata":
			g.buf.WriteFormat(`
	if _%s := r.MultipartForm.File["%s"]; len(_%s) != 0 {
		%s, err = _%s[0].Open()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
`, name, name, name, name, name)
		case "file":
			g.buf.AddImport("", "io")
			g.buf.AddImport("", "bytes")
			g.buf.WriteFormat(`
	_%s := bytes.NewBuffer(nil)
	_, err = io.Copy(_%s, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	%s = _%s
`, name, name, name, name)
		case "image":
			g.buf.AddImport("", "image")
			g.buf.AddImport("_", "image/jpeg")
			g.buf.AddImport("_", "image/png")
			g.buf.WriteFormat(`
	%s, _, err = image.Decode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`, name)

		}
	case "cookie":
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	var cookie *http.Cookie
	cookie, err = r.Cookie("%s")
	if err == nil {`, req.Name)
		g.Convert(`cookie.Value`, name, req.Type)
		g.buf.WriteFormat(`
	}
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`)

	case "query":
		g.buf.WriteFormat(`
	var _%s = r.URL.Query().Get("%s")
`, name, req.Name)
		g.Convert("_"+name, name, req.Type)

	case "header":
		g.buf.WriteFormat(`
	var _%s = r.Header.Get("%s")
`, name, req.Name)
		g.Convert("_"+name, name, req.Type)

	case "path":
		g.buf.WriteFormat(`
	var _%s = mux.Vars(r)["%s"]
`, name, req.Name)
		g.Convert("_"+name, name, req.Type)

	default:
		return fmt.Errorf("undefine in %s", req.In)
	}
	return nil
}
