package route

import (
	"fmt"
	"sort"

	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateRequestsVar(reqs []*spec.Request) error {

	for _, req := range reqs {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if req.Type == nil {
			continue
		}

		g.buf.WriteFormat("var %s ", g.getVarName(req.Name, req.Type))
		g.Types(req.Type)
		g.buf.WriteString("\n")
	}
	return nil
}

func (g *GenRoute) generateRequest(req *spec.Request, errName string) error {
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}

	vname := g.getVarName(req.Name, req.Type)
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

			if resp.Type.Ref != req.Type.Ref {
				continue
			}

			midds = append(midds, midd)
		}
		switch len(midds) {
		default:
			g.buf.WriteFormat(`
// Permission middleware undefined %s.
`, req.Name)
		case 1:
			midd := midds[0]
			name := g.getMiddlewareFunctionName(midd)
			g.buf.WriteFormat(`
// Permission middlewares call %s.
%s, %s = %s(`, midd.Name, vname, errName, name)
			if midd.Type != nil {
				g.buf.WriteString(`s, `)
			}
			g.buf.WriteFormat(`w, r)
if %s != nil {
	return
}
`, errName)
		}
	case "security":
		secus := []*spec.Security{}
		secuKey := make([]string, 0, len(g.api.Securitys))
		for k := range g.api.Securitys {
			secuKey = append(secuKey, k)
		}
		sort.Strings(secuKey)
		for _, k := range secuKey {
			secu := g.api.Securitys[k]
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
		http.Error(w, %s.Error(), 401)
`, errName)
		default:
			g.buf.WriteFormat(`
		// Permission verification
`)
			for _, secu := range secus {
				name := g.getSecurityFunctionName(secu)

				switch secu.Schema {
				case "basic":
					g.buf.AddImport("", "strings")
					g.buf.WriteFormat(`if strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") { // Call %s.
		%s, %s = %s(`, secu.Name, vname, errName, name)

				case "bearer":
					g.buf.AddImport("", "strings")
					g.buf.WriteFormat(`if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") { // Call %s.
		%s, %s = %s(`, secu.Name, vname, errName, name)

				case "apiKey":
					req := secu.Requests[0]
					if req.Ref != "" {
						req = g.api.Requests[req.Ref]
					}

					switch req.In {
					default:
					case "header":
						g.buf.WriteFormat(`if r.Header.Get("%s") != "" { // Call %s.
		%s, %s = %s(`, req.Name, secu.Name, vname, errName, name)
					case "query":
						g.buf.WriteFormat(`if r.URL.Query().Get("%s") != "" { // Call %s.
		%s, %s = %s(`, req.Name, secu.Name, vname, errName, name)
					case "cookie":
						g.buf.WriteFormat(`if cookie, err := r.Cookie("%s"); err == nil && cookie.Value != "" { // Call %s.
		%s, %s = %s(`, req.Name, secu.Name, vname, errName, name)
					}

				}
				if secu.Type != nil {
					g.buf.WriteString(`s, `)
				}
				g.buf.WriteString(`w, r)
	} else `)
			}
			g.buf.WriteFormat(` {
		http.Error(w, %s.Error(), 401)
	}
	if %s != nil {
		return
	}
`, errName, errName)
		}
	default:
		g.buf.WriteFormat(`
// Parsing %s.
%s, %s = %s(w, r)`, req.Name, vname, errName, g.getRequestFunctionName(req))
		g.buf.WriteFormat(`
if %s != nil {
	return
}
`, errName)
	}

	return nil
}

func (g *GenRoute) generateRequestFunction(req *spec.Request) error {
	g.buf.AddImport("", "net/http")

	name := g.getRequestFunctionName(req)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	g.buf.WriteFormat(`
// %s Parsing the %s for of %s
func %s(w http.ResponseWriter, r *http.Request) (%s `, name, req.In, req.Name, name, g.getVarName(req.Name, req.Type))
	g.Types(req.Type)
	g.buf.WriteString(`, err error) {
`)
	err := g.generateRequestVar(req)
	if err != nil {
		return err
	}
	g.buf.WriteString(`

	return
}`)
	return nil
}

func (g *GenRoute) generateRequestVar(req *spec.Request) error {

	name := g.getVarName(req.Name, req.Type)
	switch req.In {
	case "body":
		g.buf.WriteFormat(`
	defer r.Body.Close()
`)
		switch req.Content {
		case "json":
			g.buf.AddImport("", "io/ioutil")
			g.buf.AddImport("", "encoding/json")
			g.buf.WriteFormat(`
	var _%s []byte
	_%s, err = ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(_%s, &%s)
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
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = xml.Unmarshal(_%s, &%s)
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
			g.buf.AddImport("", "strings")
			g.buf.WriteFormat(`
	body := r.Body
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if r.MultipartForm == nil {
			err = r.ParseMultipartForm(10<<20)
			if err != nil {
				return
			}
		}
		file := r.MultipartForm.File["%s"]
		if len(file) != 0 {
			body, err = file[0].Open()
			if err != nil {
				return
			}
		}
	}

	_%s := bytes.NewBuffer(nil)
	_, err = io.Copy(_%s, body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	%s = _%s
`, req.Name, name, name, name, name)
		case "image":
			g.buf.AddImport("", "image")
			g.buf.AddImport("_", "image/jpeg")
			g.buf.AddImport("_", "image/png")
			g.buf.AddImport("", "strings")
			g.buf.WriteFormat(`
	body := r.Body
	defer r.Body.Close()
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if r.MultipartForm == nil {
			err = r.ParseMultipartForm(10<<20)
			if err != nil {
				return
			}
		}
		file := r.MultipartForm.File["%s"]
		if len(file) != 0 {
			body, err = file[0].Open()
			if err != nil {
				return
			}
		}
	}

	%s, _, err = image.Decode(body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

`, req.Name, name)

		}
	case "cookie":
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	var cookie *http.Cookie
	cookie, err = r.Cookie("%s")`, req.Name)
		g.buf.WriteFormat(`
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`)

		g.GenModel.ConvertTo(`cookie.Value`, name, req.Type)
		g.buf.WriteFormat(`
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`)
	case "query":
		varName := g.getVarName("raw_"+name, req.Type)
		g.buf.WriteFormat(`
	var %s = r.URL.Query()["%s"]
`, varName, req.Name)
		g.GenModel.ConvertToMulti(varName, name, req.Type)
		g.buf.WriteFormat(`
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`)
	case "header":
		varName := g.getVarName("raw_"+name, req.Type)
		g.buf.WriteFormat(`
	var %s = r.Header.Get("%s")
`, varName, req.Name)
		g.GenModel.ConvertTo(varName, name, req.Type)
		g.buf.WriteFormat(`
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`)
	case "path":
		varName := g.getVarName("raw_"+name, req.Type)
		g.buf.WriteFormat(`
	var %s = mux.Vars(r)["%s"]
`, varName, req.Name)
		g.GenModel.ConvertTo(varName, name, req.Type)
		g.buf.WriteFormat(`
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`)
	default:
		return fmt.Errorf("undefine in %s", req.In)
	}

	return nil
}
