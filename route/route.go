package route

import (
	"sort"
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
		api: api,
		buf: buf,
	}
}

func (g *GenRoute) Generate(pkg, outpkg, funcName string) (*srcgen.File, error) {
	if pkg == "" {
		pkg = g.api.Package
	}
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
}
`)
			}
			t = typ
			name := GetRouteName(t.Name)
			g.buf.WriteFormat(`
// %s is routing for %s
func %s(router *mux.Router,%s *%s) {
`, name, t.Name, name, GetGlobalVarName(t.Name), t.Name)
		}
		err = g.GenerateRoute(v)
		if err != nil {
			return err
		}
	}
	if t != nil {
		g.buf.WriteString(`
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
		default:
			err = g.GenerateRequestFunction(v)
			if err != nil {
				return err
			}
		}
	}

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

	name := GetGlobalVarName(typ.Name)
	g.buf.WriteFormat("// %s Define the method scope\n", typ.Name)
	g.buf.WriteFormat("var %s %s\n", name, typ.Name)
	g.buf.WriteFormat("%s(router, &%s)", GetRouteName(typ.Name), name)
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
