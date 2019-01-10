package route

import (
	"fmt"
	"unsafe"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/namecase"
)

func (g *GenRoute) GetVarName(name string) string {
	if name == "err" {
		return name

	}
	return "_" + namecase.ToCamel(fmt.Sprintf("var_%s", name))
}

func (g *GenRoute) GetRouteName(typ *spec.Type) string {
	name := namecase.ToPascal(fmt.Sprintf("route_%s", typ.Name))
	addr := unsafe.Pointer(typ)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) GetOperationFunctionName(oper *spec.Operation) string {
	name := "_" + namecase.ToCamel(fmt.Sprintf("operation_%s_%s", oper.Method, oper.Path))
	addr := unsafe.Pointer(oper)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) GetRequestFunctionName(req *spec.Request) string {
	name := "_" + namecase.ToCamel(fmt.Sprintf("request_%s_%s", req.In, req.Name))
	addr := unsafe.Pointer(req)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) GetSecurityFunctionName(secu *spec.Security) string {
	name := "_" + namecase.ToCamel(fmt.Sprintf("security_%s_%s", secu.Schema, secu.Name))
	addr := unsafe.Pointer(secu)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) GetMiddlewareFunctionName(midd *spec.Middleware) string {
	name := "_" + namecase.ToCamel(fmt.Sprintf("middleware_%s_%s", midd.Schema, midd.Name))
	addr := unsafe.Pointer(midd)
	return g.named.GetName(name, addr)
}
