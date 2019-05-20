package route

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/wzshiming/gen/model"
	"github.com/wzshiming/gen/named"
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/srcgen"
	"github.com/wzshiming/openapi/util"
)

// GenClient is the generating generating
type GenRoute struct {
	api          *spec.API
	apiInterface interface{}
	buf          *srcgen.File
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
	g.GenModel = *model.NewGenModel(g.api, g.buf, []string{outpkg})
	err := g.generateRoutes(funcName)
	if err != nil {
		return nil, err
	}
	if g.apiInterface != nil {
		err := g.generateOpenAPI()
		if err != nil {
			return nil, err
		}
	}
	return g.buf, nil
}

func (g *GenRoute) generateOpenAPI() (err error) {
	dj, err := json.Marshal(g.apiInterface)
	if err != nil {
		return err
	}
	dy, err := util.JSON2YAML(dj)
	if err != nil {
		return err
	}
	oay := strings.Replace(string(dy), "`", "`+\"`\"+`", -1)
	oaj := strings.Replace(string(dj), "`", "`+\"`\"+`", -1)
	g.buf.WriteFormat("var OpenAPI4YAML=[]byte(`%s`)\n", oay)
	g.buf.WriteFormat("var OpenAPI4JSON=[]byte(`%s`)\n", oaj)

	g.buf.AddImport("", "github.com/wzshiming/openapi/ui")
	g.buf.AddImport("", "github.com/wzshiming/openapi/ui/swaggerui")
	g.buf.AddImport("", "github.com/wzshiming/openapi/ui/swaggereditor")
	g.buf.AddImport("", "github.com/wzshiming/openapi/ui/redoc")
	g.buf.AddImport("", "github.com/gorilla/mux")

	g.buf.WriteString(`
// RouteOpenAPI
func RouteOpenAPI(router *mux.Router) *mux.Router {
	openapi := map[string][]byte {
		"openapi.json": OpenAPI4JSON,
		"openapi.yml": OpenAPI4YAML,
		"openapi.yaml": OpenAPI4YAML,
	}
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", ui.HandleWithFiles(openapi, swaggerui.Asset)))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui", ui.HandleWithFiles(openapi, swaggerui.Asset)))
	router.PathPrefix("/swaggereditor/").Handler(http.StripPrefix("/swaggereditor", ui.HandleWithFiles(openapi, swaggereditor.Asset)))
	router.PathPrefix("/redoc/").Handler(http.StripPrefix("/redoc", ui.HandleWithFiles(openapi, redoc.Asset)))
	return router
}
`)
	return nil
}

func (g *GenRoute) WithOpenAPI(api interface{}) *GenRoute {
	g.apiInterface = api
	return g
}

func (g *GenRoute) generateSubRoutes(basePath string, op []*spec.Operation) (err error) {

	sort.Slice(op, func(i, j int) bool {
		io := op[i].Path
		jo := op[j].Path
		if io == jo {
			return true
		}
		ip := strings.Count(io, "{")
		jp := strings.Count(jo, "{")
		if ip == jp {
			return io > jo
		}
		return ip < jp
	})

	name := g.getVarName("route_"+basePath, op[0].Type)

	g.buf.WriteFormat(`
		%s := router.PathPrefix("%s").Subrouter()
		if len(fs) != 0 {
			%s.Use(fs...)
		}
`, name, basePath, name)

	for _, v := range op {
		typ := v.Type
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}

		err = g.generateRoute(name, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GenRoute) generateRouteFunc(typ *spec.Type, op []*spec.Operation) (err error) {
	name := g.getRouteName(typ)
	g.buf.WriteFormat(`
	// %s is routing for %s
	func %s(router *mux.Router, %s `, name, typ.Name, name, g.getVarName("", typ))
	g.PtrTypes(typ)
	g.buf.WriteFormat(`, fs ...mux.MiddlewareFunc) *mux.Router {
	if router == nil {
		router = mux.NewRouter()
	}
`)

	group := map[string][]*spec.Operation{}
	sortGroup := []string{}
	for _, v := range op {
		if group[v.BasePath] == nil {
			sortGroup = append(sortGroup, v.BasePath)
		}
		group[v.BasePath] = append(group[v.BasePath], v)
	}

	sort.Slice(sortGroup, func(i, j int) bool {
		io := sortGroup[i]
		jo := sortGroup[j]
		if io == jo {
			return true
		}
		ip := strings.Count(io, "{")
		jp := strings.Count(jo, "{")
		if ip == jp {
			return io > jo
		}
		return ip < jp
	})

	for _, basePath := range sortGroup {
		op := group[basePath]
		err := g.generateSubRoutes(basePath, op)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(`
	return router
}
`)
	return nil
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

	if g.apiInterface != nil {
		g.buf.WriteString(`
	router = RouteOpenAPI(router)
		`)
	}
	g.buf.WriteString(`
	return router
}
`)

	group := map[*spec.Type][]*spec.Operation{}
	sortGroup := []*spec.Type{}
	for _, v := range g.api.Operations {
		if v.Type == nil {
			continue
		}
		typ := v.Type
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		if group[typ] == nil {
			sortGroup = append(sortGroup, typ)
		}
		group[typ] = append(group[typ], v)
	}

	for _, typ := range sortGroup {
		op := group[typ]
		err := g.generateRouteFunc(typ, op)
		if err != nil {
			return err
		}
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
		case "wrapping":
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
		wrapKey := make([]string, 0, len(g.api.Wrappings))
		for k := range g.api.Wrappings {
			wrapKey = append(wrapKey, k)
		}
		sort.Strings(wrapKey)

		for _, k := range wrapKey {
			v := g.api.Wrappings[k]
			err = g.generateWrappingFunction(v)
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
			secu := g.api.Securitys[k]
			err = g.generateSecurityFunction(secu)
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
%s(router, `, g.getRouteName(typ))
	if g.Ptr(typ) {
		g.buf.WriteString("&")
	}
	g.buf.WriteFormat(`%s)
`, name)
	m[typ.Name] = true
	return
}

func (g *GenRoute) generateRoute(subrouter string, oper *spec.Operation) (err error) {
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
	%s.Methods("%s").Path("%s").Handler(_%s)
`, subrouter, strings.Join(methods, `", "`), path, name)
	return
}
