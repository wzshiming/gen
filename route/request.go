package route

import (
	"fmt"

	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) GenerateRequestFunction(req *spec.Request) error {
	g.buf.AddImport("", "net/http")

	name := g.GetRequestFunctionName(req)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	g.buf.WriteFormat(`
// %s Parsing the %s for of %s
func %s(w http.ResponseWriter, r *http.Request) (%s `, name, req.In, req.Name, name, req.Name)
	g.Types(req.Type)
	g.buf.WriteString(`,err error) {
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

	switch req.In {
	case "body":
		g.buf.AddImport("", "io/ioutil")
		g.buf.AddImport("", "encoding/json")
		g.buf.WriteFormat(`
	var _body []byte
	_body, err = ioutil.ReadAll(r.Body)
	if err == nil {
		r.Body.Close()
		err = json.Unmarshal(_body, &%s)
	}
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
`, req.Name)

	case "cookie":
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	var _cookie *http.Cookie
	_cookie, err = r.Cookie("%s")
	if err == nil {`, req.Name)
		g.Convert(`_cookie.Value`, req.Name, req.Type)
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
`, req.Name, req.Name)
		g.Convert("_"+req.Name, req.Name, req.Type)

	case "header":
		g.buf.WriteFormat(`
	var _%s = r.Header.Get("%s")
`, req.Name, req.Name)
		g.Convert("_"+req.Name, req.Name, req.Type)

	case "path":
		g.buf.WriteFormat(`
	var _%s = mux.Vars(r)["%s"]
`, req.Name, req.Name)
		g.Convert("_"+req.Name, req.Name, req.Type)

	default:
		return fmt.Errorf("undefine in %s", req.In)
	}
	return nil
}
