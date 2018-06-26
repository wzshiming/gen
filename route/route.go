package route

import (
	"strings"

	"github.com/wzshiming/ffmt"
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
	defer g.buf.WriteString(`
	return router
}
`)

	for _, v := range g.api.Operations {
		g.GenerateRoute(v)
	}
	return
}

func (g *GenRoute) GenerateRoute(oper *spec.Operation) (err error) {
	g.buf.WriteFormat(`router.Path("%s").`, oper.Path)
	g.buf.WriteFormat(`Methods("%s").`, strings.ToUpper(oper.Method))
	for _, req := range oper.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		switch req.In {
		case "query":
			g.buf.WriteFormat(`Queries("%s").`, req.Name)
		case "header":
			g.buf.WriteFormat(`Headers("%s").`, req.Name)
		}
	}

	g.buf.AddImport("", "net/http")
	g.buf.WriteString(`
HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_vars := mux.Vars(r)
`)
	defer g.buf.WriteString(`
})
`)

	for _, req := range oper.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		g.Convert(`_vars["`+req.Name+`"]`, "_"+req.Name, req.Type)
		g.buf.WriteString("\n")
	}
	g.buf.WriteString("\n")

	for i, resp := range oper.Responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		g.buf.WriteString("_" + resp.Name)
	}

	g.buf.WriteString(":=")
	g.buf.WriteString(oper.Name)

	g.buf.WriteString("(")
	for i, req := range oper.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if req.Type.Kind == spec.Ptr {
			g.buf.WriteString("&")
		}
		g.buf.WriteString("_" + req.Name)
	}
	g.buf.WriteString(")")

	for _, resp := range oper.Responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}

		g.MarshalContentBody(resp)
	}

	return
}

func (g *GenRoute) UnmarshalContentBody(req *spec.Request) error {
	g.buf.AddImport("", "io/ioutil")
	g.buf.WriteFormat(`func()(_b %s){
	data, err := ioutil.ReadAll(_body)
	if err != nil {
		return
	}
	_body.Close()
	`, g.Types(req.Type))
	switch req.Content {
	case "json":
		g.buf.WriteString(`json.Unmarshal(data,&_b)`)
	case "xml":
		g.buf.WriteString(`xml.Unmarshal(data,&_b)`)
	}
	g.buf.WriteString(`
	return
}()`)
	return nil
}

func (g *GenRoute) MarshalContentBody(resp *spec.Response) error {
	g.buf.WriteFormat(`
	if %s != nil {`, resp.Name)

	contentType := ""
	switch resp.Content {
	case "json":
		contentType = "application/json; charset=utf-8"
		g.buf.WriteString(`
	data, err:=json.Marshal(_b)
	if err!=nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}`)
		g.buf.AddImport("", "encoding/json")
	case "xml":
		contentType = "application/xml; charset=utf-8"
		g.buf.WriteString(`
	data,err:=xml.Marshal(_b)
	if err!=nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}`)
		g.buf.AddImport("", "encoding/xml")
	default:
		contentType = "text/plain; charset=utf-8"
	}
	g.buf.WriteFormat(`
	w.Header().Set("Content-Type","%s")
	w.WriteHeader(%s)
	w.Write([]byte(err.Error()))
	return
}
`, contentType, resp.Code)
	return nil
}

func (g *GenRoute) Convert(in, out string, typ *spec.Type) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	switch typ.Kind {
	case spec.Ptr:
		return g.Convert(in, out, typ.Elem)
	case spec.String:
		g.buf.WriteFormat(`%s := %s`, out, in)
	case spec.Int64:
		g.buf.AddImport("", "strconv")
		g.buf.WriteFormat("%s, _ := strconv.ParseInt(%s,0,0)\n", out, in)
	case spec.Int8, spec.Int16, spec.Int32, spec.Int:
		g.buf.AddImport("", "strconv")
		g.buf.WriteFormat("_%s, _ := strconv.ParseInt(%s,0,0)\n", out, in)
		g.buf.WriteFormat("%s := %s(_%s)", out, typ.Kind, out)
	case spec.Uint64:
		g.buf.AddImport("", "strconv")
		g.buf.WriteFormat("%s, _ := strconv.ParseUint(%s,0,0)\n", out, in)
	case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint:
		g.buf.AddImport("", "strconv")
		g.buf.WriteFormat("_%s, _ := strconv.ParseUint(%s,0,0)\n", out, in)
		g.buf.WriteFormat("%s := %s(_%s)\n", out, typ.Kind, out)
	case spec.Slice:
		if typ.Elem.Kind == spec.Byte {
			g.buf.WriteFormat("%s := []byte(%s)\n", out, in)
		}
	default:
		ffmt.P(typ)
		//		g.buf.WriteString(typ.Type)
	}
	return nil
}
