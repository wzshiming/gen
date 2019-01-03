package route

import (
	"sort"
	"strings"

	"github.com/wzshiming/gen/model"
	"github.com/wzshiming/gen/named"
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/srcgen"
)

// GenClient is the generating generating
type GenRoute struct {
	api *spec.API
	buf *srcgen.File
	model.GenModel
	only  map[string]bool
	named *named.Named
}

func NewGenRoute(api *spec.API) *GenRoute {
	buf := &srcgen.File{}
	return &GenRoute{
		api:   api,
		buf:   buf,
		only:  map[string]bool{},
		named: named.NewNamed(),
	}
}

func (g *GenRoute) Generate(pkg, outpkg, funcName string) (*srcgen.File, error) {
	g.buf.WithPackname(pkg)
	g.GenModel = *model.NewGenModel(g.api, g.buf, outpkg)
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
	g.buf.WriteFormat(`
// %s is all routing for package
// generated do not edit.
func %s() http.Handler {
	router := mux.NewRouter()
`, funcName, funcName)

	for _, v := range g.api.Operations {
		err = g.GenerateRouteTypes(v, m)
		if err != nil {
			return err
		}
	}

	for _, v := range g.api.Operations {
		if v.Type != nil {
			continue
		}
		err = g.GenerateRoute(v)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(`
	return router
}
`)

	var t *spec.Type
	for _, v := range g.api.Operations {
		typ := v.Type
		if typ == nil {
			continue
		}
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}

		if typ != t {
			if t != nil {
				g.buf.WriteString(`
	return router
}
`)
			}
			t = typ
			name := g.GetRouteName(t)
			g.buf.WriteFormat(`
// %s is routing for %s
func %s(router *mux.Router, %s *`, name, t.Name, name, g.GetVarName(t))
			g.Types(v.Type)
			g.buf.WriteFormat(`) *mux.Router {
	if router == nil {
		router = mux.NewRouter()
	}
`)
		}
		err = g.GenerateRoute(v)
		if err != nil {
			return err
		}
	}
	if t != nil {
		g.buf.WriteString(`
	return router
}
`)
	}

	reqKey := make([]string, 0, len(g.api.Requests))
	for k := range g.api.Requests {
		reqKey = append(reqKey, k)
	}
	sort.Strings(reqKey)
	for _, k := range reqKey {
		v := g.api.Requests[k]
		switch v.In {
		case "security":
		case "middleware":
		case "none":
		default:
			err = g.GenerateRequestFunction(v)
			if err != nil {
				return err
			}
		}
	}

	{
		middKey := make([]string, 0, len(g.api.Middlewares))
		for k := range g.api.Middlewares {
			middKey = append(middKey, k)
		}
		sort.Strings(middKey)

		for _, k := range middKey {
			v := g.api.Middlewares[k]
			err = g.GenerateMiddlewareFunction(v)
			if err != nil {
				return err
			}
		}
	}

	{
		secuKey := make([]string, 0, len(g.api.Securitys))
		for k := range g.api.Securitys {
			secuKey = append(secuKey, k)
		}

		sort.Strings(secuKey)
		for _, k := range secuKey {
			v := g.api.Securitys[k]
			err = g.GenerateSecurityFunction(v)
			if err != nil {
				return err
			}
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

	name := g.GetVarName(typ)
	g.buf.WriteFormat(`
// %s Define the method scope
var %s `, typ.Name, name)
	g.Types(oper.Type)
	g.buf.WriteFormat(`
%s(router, &%s)
`, g.GetRouteName(typ), name)
	m[typ.Name] = true
	return
}

func (g *GenRoute) GenerateRoute(oper *spec.Operation) (err error) {
	name := g.GetOperationFunctionName(oper)

	methods := strings.Split(oper.Method, ",")
	for i := range methods {
		methods[i] = `"` + strings.ToUpper(methods[i]) + `"`
	}

	g.buf.WriteFormat(`
	// Registered routing %s %s
	router.Methods(%s).Path("%s").
`, strings.ToUpper(oper.Method), oper.Path, strings.Join(methods, ", "), oper.Path)

	if oper.Type != nil {
		typ := oper.Type
		if typ.Ref != "" {
			typ = g.api.Types[oper.Type.Ref]
		}
		typname := g.GetVarName(typ)
		g.buf.WriteFormat(`
		HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			%s(%s, w, r)
		})
`, name, typname)
	} else {
		g.buf.WriteFormat(`
		HandlerFunc(%s)
`, name)
	}
	return
}
