package route

import (
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

func (g *GenRoute) Generate() ([]byte, error) {
	g.buf.WithPackname(g.api.Package)
	err := g.GenerateRoutes()
	if err != nil {
		return nil, err
	}

	return g.buf.Bytes(), nil
}

func (g *GenRoute) GenerateRoutes() (err error) {
	g.buf.AddImport("", "github.com/gorilla/mux")

	g.buf.WriteString(`
func Router() *mux.Router {
	router := mux.NewRouter()
`)
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

func (g *GenRoute) GenerateRoute(oper *spec.Operation) (err error) {
	g.buf.WriteFormat(`
	// Registered routing %s %s
	router.Path("%s").
		Methods("%s").
		HandlerFunc(_%s)
`, strings.ToUpper(oper.Method), oper.Path, oper.Path, strings.ToUpper(oper.Method), oper.Name)
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
func _%s(w http.ResponseWriter, r *http.Request) {
`, oper.Name, oper.Name, oper.Name)
	defer g.buf.WriteString(`
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
		g.Convert(`cookie.Value`, "_"+req.Name, req.Type, true)
		g.buf.WriteFormat(`}
`)
	case "query":
		g.buf.WriteFormat(`
	// Parsing the query for %s.
`, req.Name)
		g.Convert(`r.URL.Query().Get("`+req.Name+`")`, "_"+req.Name, req.Type, false)
	case "header":
		g.buf.WriteFormat(`
	// Parsing the header for %s.
`, req.Name)
		g.Convert(`r.URL.Header.Get("`+req.Name+`")`, "_"+req.Name, req.Type, false)
	default:
		g.buf.WriteFormat(`
	// Parsing the %s for %s.
`, req.In, req.Name)
		g.Convert(`mux.Vars(r)["`+req.Name+`"]`, "_"+req.Name, req.Type, false)
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

	g.buf.WriteFormat(":= %s(", oper.Name)
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

func (g *GenRoute) Convert(in, out string, typ *spec.Type, defined bool) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	switch typ.Kind {
	case spec.Ptr:
		return g.Convert(in, out, typ.Elem, defined)
	case spec.String:
		if !defined {
			g.buf.WriteFormat(`%s := %s`, out, in)
		} else {
			g.buf.WriteFormat(`%s = %s`, out, in)
		}
		return nil
	case spec.Int64:
		g.buf.AddImport("", "strconv")
		if !defined {
			g.buf.WriteFormat("var %s %s\n", out, strings.ToLower(typ.Kind.String()))
		}
		g.buf.WriteFormat(`if i, err := strconv.ParseInt(%s,0,0); err == nil {
	%s = i
}
`, in, out)
		return nil
	case spec.Int8, spec.Int16, spec.Int32, spec.Int:
		g.buf.AddImport("", "strconv")
		if !defined {
			g.buf.WriteFormat("var %s %s\n", out, strings.ToLower(typ.Kind.String()))
		}
		g.buf.WriteFormat(`if i, err := strconv.ParseInt(%s,0,0); err == nil {
	%s = %s(i)
}
`, in, out, strings.ToLower(typ.Kind.String()))
		return nil
	case spec.Uint64:
		g.buf.AddImport("", "strconv")
		if !defined {
			g.buf.WriteFormat("var %s %s\n", out, strings.ToLower(typ.Kind.String()))
		}
		g.buf.WriteFormat(`if i, err := strconv.ParseUint(%s,0,0); err == nil {
	%s = i
}
`, in, out)
		return nil
	case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint:
		g.buf.AddImport("", "strconv")
		if !defined {
			g.buf.WriteFormat("var %s %s\n", out, strings.ToLower(typ.Kind.String()))
		}
		g.buf.WriteFormat(`if i, err := strconv.ParseUint(%s,0,0); err == nil {
	%s = %s(i)
}
`, in, out, strings.ToLower(typ.Kind.String()))
		return nil
	case spec.Slice:
		if typ.Elem.Kind == spec.Byte {
			g.buf.WriteFormat("%s := []byte(%s)\n", out, in)
			return nil
		}
	default:
	}

	g.buf.WriteFormat("// Conversion of string to ")
	g.Types(typ)
	g.buf.WriteString(" is not supported.")
	g.buf.WriteFormat("\nvar %s ", out)
	g.Types(typ)

	return nil
}
