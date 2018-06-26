package route

import (
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
HandlerFunc(func(_w http.ResponseWriter, _r *http.Request) {
	_vars := mux.Vars(_r)
`)
	defer g.buf.WriteString(`
})
`)

	for i, resp := range oper.Responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		g.buf.WriteString(resp.Name)
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
		switch req.In {
		case "body":
			err := g.UnmarshalContentBody(req)
			if err != nil {
				return err
			}

		case "cookie":
			// TODO
			g.buf.WriteString("nil")
		default:
			g.buf.WriteFormat(`_vars["%s"]`, req.Name)
		}
	}
	g.buf.WriteString(")")

	for _, resp := range oper.Responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}

		g.MarshalContentBody(resp)

		//g.buf.WriteString(resp.Name)
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
		_w.WriteHeader(500)
		_w.Write([]byte(err.Error()))
		return
	}`)
		g.buf.AddImport("", "encoding/json")
	case "xml":
		contentType = "application/xml; charset=utf-8"
		g.buf.WriteString(`
	data,err:=xml.Marshal(_b)
	if err!=nil {
		_w.WriteHeader(500)
		_w.Write([]byte(err.Error()))
		return
	}`)
		g.buf.AddImport("", "encoding/xml")
	default:
		contentType = "text/plain; charset=utf-8"
	}
	g.buf.WriteFormat(`
	_w.Header().Set("Content-Type","%s")
	_w.WriteHeader(%s)
	_w.Write([]byte(err.Error()))
	return
}
`, contentType, resp.Code)
	return nil
}
