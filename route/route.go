package route

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wzshiming/gen/model"
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/srcgen"
)

// GenClient is the generating generating
type GenRoute struct {
	api *spec.API
	buf *srcgen.File
	model.GenModel
}

func NewGenRoute(api *spec.API) *GenRoute {
	buf := &srcgen.File{}
	return &GenRoute{
		api:      api,
		buf:      buf,
		GenModel: *model.NewGenModel(api, buf),
	}
}

func (g *GenRoute) Generate(funcName string) (*srcgen.File, error) {
	g.buf.WithPackname(g.api.Package)
	err := g.GenerateRoutes(funcName)
	if err != nil {
		return nil, err
	}

	return g.buf, nil
}

func (g *GenRoute) GenerateRoutes(funcName string) (err error) {
	g.buf.AddImport("", "github.com/gorilla/mux")
	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
func %s() http.Handler {
	router := mux.NewRouter()

`, funcName)
	m := map[string]bool{}
	for _, v := range g.api.Operations {
		err = g.GenerateRouteTypes(v, m)
		if err != nil {
			return err
		}
	}
	for _, v := range g.api.Operations {
		err = g.GenerateRoute(v)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(`
	return router
}
`)

	for _, v := range g.api.Operations {
		err = g.GenerateRouteFunction(v)
		if err != nil {
			return err
		}
	}

	return
}

func (g *GenRoute) GenerateRouteTypes(oper *spec.Operation, m map[string]bool) (err error) {
	if oper.Type == nil {
		return
	}
	typ := oper.Type
	if typ.Ref != "" {
		typ = g.api.Types[oper.Type.Ref]
	}
	if m[typ.Name] {
		return
	}
	g.buf.WriteFormat("// %s Define the method scope\n", typ.Name)
	g.buf.WriteFormat("var _%s %s\n", typ.Name, typ.Name)
	m[typ.Name] = true
	return
}

func (g *GenRoute) GenerateRoute(oper *spec.Operation) (err error) {
	g.buf.WriteFormat(`
	// Registered routing %s %s
	router.Path("%s").
		Methods("%s").
		HandlerFunc(`, strings.ToUpper(oper.Method), oper.Path, oper.Path, strings.ToUpper(oper.Method))
	if oper.Type != nil {
		typ := oper.Type
		if typ.Ref != "" {
			typ = g.api.Types[oper.Type.Ref]
		}
		g.buf.WriteString("_")
		g.Types(oper.Type)
		g.buf.WriteString(".")
	}
	g.buf.WriteFormat(`_%s)
`, oper.Name)
	return
}

func (g *GenRoute) GenerateRouteFunc(oper *spec.Operation) (err error) {
	g.buf.WriteFormat(`HandlerFunc(_%s)`, oper.Name)
	return
}

func (g *GenRoute) GenerateRouteFunction(oper *spec.Operation) (err error) {

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
// _%s Is the route of %s 
func`, oper.Name, oper.Name)
	if oper.Type != nil {
		g.buf.WriteString("(s ")
		g.Types(oper.Type)
		g.buf.WriteString(")")
	}
	g.buf.WriteFormat(` _%s(w http.ResponseWriter, r *http.Request) {
`, oper.Name)
	defer g.buf.WriteString(`
	w.WriteHeader(204)
	w.Write(nil)
	return
}
`)

	for _, req := range oper.Requests {
		err = g.GenerateRequest(req)
		if err != nil {
			return err
		}
	}
	err = g.GenerateCall(oper)
	if err != nil {
		return err
	}
	for _, resp := range oper.Responses {
		err = g.GenerateResponse(resp)
		if err != nil {
			return err
		}
	}
	return
}

func (g *GenRoute) GenerateRequest(req *spec.Request) error {
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}
	switch req.In {
	case "body":
		g.buf.AddImport("", "io/ioutil")
		g.buf.WriteFormat(`
	// Parsing the body for %s.
	var _%s `, req.Name, req.Name)
		g.Types(req.Type)
		g.buf.WriteFormat(`
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		r.Body.Close()
		json.Unmarshal(body, &_%s)
	}
`, req.Name)
	case "cookie":
		g.buf.WriteFormat(`
	// Parsing the cookie for %s.
	var _%s `, req.Name, req.Name)
		g.Types(req.Type)
		g.buf.WriteFormat(`
	if cookie, err := r.Cookie("%s"); err == nil {`, req.Name)
		g.Convert(`cookie.Value`, "_"+req.Name, req.Type)
		g.buf.WriteFormat(`}
`)

	case "query":
		g.buf.WriteFormat(`
	// Parsing the query for %s.
	var _in_%s = r.URL.Query().Get("%s")`, req.Name, req.Name, req.Name)
		g.buf.WriteFormat(`
	var _%s `, req.Name)
		g.Types(req.Type)
		g.buf.WriteString("\n")
		g.Convert("_in_"+req.Name, "_"+req.Name, req.Type)

	case "header":
		g.buf.WriteFormat(`
	// Parsing the header for %s.
	var _in_%s = r.URL.Header.Get("%s")`, req.Name, req.Name, req.Name)
		g.buf.WriteFormat(`
	var _%s `, req.Name)
		g.Types(req.Type)
		g.buf.WriteString("\n")
		g.Convert("_in_"+req.Name, "_"+req.Name, req.Type)

	case "path":
		g.buf.WriteFormat(`
	// Parsing the path for %s.
	var _in_%s = mux.Vars(r)["%s"]`, req.Name, req.Name, req.Name)
		g.buf.WriteFormat(`
	var _%s `, req.Name)
		g.Types(req.Type)
		g.buf.WriteString("\n")
		g.Convert("_in_"+req.Name, "_"+req.Name, req.Type)

	default:
		return fmt.Errorf("undefine in %s", req.In)
	}
	g.buf.WriteString("\n")
	return nil
}

func (g *GenRoute) GenerateCall(oper *spec.Operation) error {
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
	}
	g.buf.WriteFormat("%s(", oper.Name)
	for i, req := range oper.Requests {
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

func (g *GenRoute) GenerateResponse(resp *spec.Response) error {
	if resp.Ref != "" {
		resp = g.api.Responses[resp.Ref]
	}
	contentType := ""
	errResp := func() {
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
`)
	}

	text := ""
	if i, err := strconv.Atoi(resp.Code); err == nil {
		text = http.StatusText(i)
	}
	g.buf.WriteFormat(`
	// Response code %s %s for %s.
	if _%s != `, resp.Code, text, resp.Name, resp.Name)
	g.TypesZero(resp.Type)
	g.buf.WriteFormat(`{`)
	switch resp.Content {
	case "json":
		g.buf.AddImport("", "encoding/json")
		contentType = "application/json; charset=utf-8"
		g.buf.WriteFormat(`
	data, err := json.Marshal(_%s)`, resp.Name)
		errResp()
	case "xml":
		g.buf.AddImport("", "encoding/xml")
		contentType = "application/xml; charset=utf-8"
		g.buf.WriteFormat(`
	data, err := xml.Marshal(_%s)`, resp.Name)
		errResp()
	case "error":
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	http.Error(w, _%s.Error(), %s)
	return
}
`, resp.Name, resp.Code)
		return nil
	default:
		contentType = "text/plain; charset=utf-8"
	}

	g.buf.WriteFormat(`
		w.Header().Set("Content-Type","%s")
		w.WriteHeader(%s)
		w.Write(data)
		return
	}
`, contentType, resp.Code)
	return nil
}

func (g *GenRoute) convertString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = %s`, out, in)
	return nil
}

func (g *GenRoute) convertPrtString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = &%s`, out, in)
	return nil
}

func (g *GenRoute) convertInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseInt(%s,0,0); err == nil {
	%s = %s(i)
}
`, in, out, name)
	return nil
}

func (g *GenRoute) convertPrtInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseInt(%s,0,0); err == nil {
	_i := %s(i)
	%s = &_i
}
`, in, name, out)
	return nil
}

func (g *GenRoute) convertUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseUint(%s,0,0); err == nil {
	%s = %s(i)
}
`, in, out, name)
	return nil
}

func (g *GenRoute) convertPrtUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseUint(%s,0,0); err == nil {
	_i := %s(i)
	%s = &_i
}
`, in, name, out)
	return nil
}

func (g *GenRoute) Convert(in, out string, typ *spec.Type) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	switch typ.Kind {
	case spec.Ptr:
		typ = typ.Elem
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		switch typ.Kind {
		case spec.String:
			return g.convertPrtString(in, out, typ)
		case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
			return g.convertPrtInt64(in, out, typ)
		case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
			return g.convertPrtUint64(in, out, typ)
		default:
		}
	default:
		switch typ.Kind {
		case spec.String:
			return g.convertString(in, out, typ)
		case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
			return g.convertInt64(in, out, typ)
		case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
			return g.convertUint64(in, out, typ)
		case spec.Slice:
			if typ.Elem.Kind == spec.Byte {
				g.buf.WriteFormat("%s := []byte(%s)\n", out, in)
				return nil
			}
		}
	}

	g.buf.WriteFormat("// Conversion of string to ")
	g.Types(typ)
	g.buf.WriteString(" is not supported.")
	g.buf.WriteFormat("\nvar %s ", out)
	g.Types(typ)

	return nil
}
