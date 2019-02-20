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
		named: named.NewNamed("_"),
	}
}

func (g *GenRoute) Generate(pkg, outpkg, funcName string) (*srcgen.File, error) {
	g.buf.WithPackname(pkg)
	g.GenModel = *model.NewGenModel(g.api, g.buf, outpkg)
	err := g.generateRoutes(funcName)
	if err != nil {
		return nil, err
	}

	return g.buf, nil
}

func (g *GenRoute) generateRoutes(funcName string) (err error) {
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
		err = g.generateRouteTypes(v, m)
		if err != nil {
			return err
		}
	}

	for _, v := range g.api.Operations {
		if v.Type != nil {
			continue
		}
		err = g.generateRoute(v)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(`
	return router
}
`)

	opers := map[string][]*spec.Operation{}
	for _, v := range g.api.Operations {
		if v.Type == nil {
			continue
		}
		opers[v.BasePath] = append(opers[v.BasePath], v)
	}

	basepaths := []string{}
	for k, _ := range opers {
		basepaths = append(basepaths, k)
	}
	sort.Strings(basepaths)

	for _, p := range basepaths {
		op := opers[p]

		sort.Slice(op, func(i, j int) bool {
			io := op[i].Path
			jo := op[j].Path
			ip := strings.Count(io, "{")
			jp := strings.Count(jo, "{")
			if ip == jp {
				return io > jo
			}
			return ip < jp
		})

		for i, v := range op {
			typ := v.Type
			if typ.Ref != "" {
				typ = g.api.Types[typ.Ref]
			}

			if i == 0 {
				name := g.getRouteName(typ)
				g.buf.WriteFormat(`
// %s is routing for %s
func %s(router *mux.Router, %s *`, name, typ.Name, name, g.getVarName("", typ))
				g.Types(v.Type)
				g.buf.WriteFormat(`, fs ...mux.MiddlewareFunc) *mux.Router {
	if router == nil {
		router = mux.NewRouter()
	}
	subrouter := router.PathPrefix("%s").Subrouter()
	if len(fs) != 0 {
		subrouter.Use(fs...)
	}
`, v.BasePath)
			}

			err = g.generateRoute(v)
			if err != nil {
				return err
			}
		}
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
			err = g.generateRequestFunction(v)
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
			err = g.generateMiddlewareFunction(v)
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
			err = g.generateSecurityFunction(v)
			if err != nil {
				return err
			}
		}
	}

	for _, v := range g.api.Operations {
		err = g.generateOperationFunction(v)
		if err != nil {
			return err
		}
	}

	return
}

func (g *GenRoute) generateRouteTypes(oper *spec.Operation, m map[string]bool) (err error) {
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

	name := g.getVarName("", typ)
	g.buf.WriteFormat(`
// %s Define the method scope
var %s `, typ.Name, name)
	g.Types(oper.Type)
	g.buf.WriteFormat(`
%s(router, &%s)
`, g.getRouteName(typ), name)
	m[typ.Name] = true
	return
}

func (g *GenRoute) generateRoute(oper *spec.Operation) (err error) {
	name := g.getOperationFunctionName(oper)

	methods := strings.Split(oper.Method, ",")
	for i := range methods {
		methods[i] = strings.ToUpper(methods[i])
	}

	path := oper.Path[len(oper.BasePath):]
	g.buf.WriteFormat(`
	// Registered routing %s %s`, strings.ToUpper(oper.Method), oper.Path)
	if oper.Type != nil {
		typ := oper.Type
		if typ.Ref != "" {
			typ = g.api.Types[oper.Type.Ref]
		}
		typname := g.getVarName("", typ)
		g.buf.WriteFormat(`
	var _%s http.Handler
	_%s = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		%s(%s, w, r)
	})`, name, name, name, typname)
	} else {
		g.buf.WriteFormat(`
	_%s := http.HandlerFunc(%s)`, name, name)
	}
	g.buf.WriteFormat(`
	subrouter.Methods("%s").Path("%s").Handler(_%s)
`, strings.Join(methods, `", "`), path, name)
	return
}
