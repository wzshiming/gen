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

	m := map[string]bool{}
	for _, v := range g.api.Operations {
		err = g.GenerateRouteTypes(v, m)
		if err != nil {
			return err
		}
	}

	g.buf.WriteFormat(`
// %s is generated do not edit.
func %s() http.Handler {
	router := mux.NewRouter()

`, funcName, funcName)
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

	for _, v := range g.api.Requests {
		switch v.In {
		case "security":
		default:
			err = g.GenerateRequestFunction(v)
			if err != nil {
				return err
			}
		}
	}

	for _, v := range g.api.Securitys {
		err = g.GenerateSecurityFunction(v)
		if err != nil {
			return err
		}
	}

	for _, v := range g.api.Operations {
		err = g.GenerateOperationFunction(v)
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
	g.buf.WriteFormat("var %s %s\n", GetGlobalVarName(typ.Name), typ.Name)
	m[typ.Name] = true
	return
}

func (g *GenRoute) GenerateRoute(oper *spec.Operation) (err error) {
	name := GetOperationFunctionName(oper.Name)
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
		g.buf.WriteFormat("%s.", GetGlobalVarName(typ.Name))
	}
	g.buf.WriteFormat(`%s)
`, name)
	return
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
