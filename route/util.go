package route

import (
	"fmt"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/namecase"
)

func (g *GenRoute) GetVarName(typ *spec.Type) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("var_%s", typ.Name))
}

func (g *GenRoute) GetRouteName(typ *spec.Type) string {
	return namecase.ToPascal(fmt.Sprintf("route_%s", typ.Name))
}

func (g *GenRoute) GetOperationFunctionName(oper *spec.Operation) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("operation_%s_%s_%s", oper.Method, oper.PkgPath, oper.Name))
}

func (g *GenRoute) GetRequestFunctionName(req *spec.Request) string {
	typ := req.Type
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	return "_" + namecase.ToCamel(fmt.Sprintf("request_%s_%s_%s_%s_%s", req.In, req.Name, typ.PkgPath, typ.Name, typ.Kind.String()))

}

func (g *GenRoute) GetSecurityFunctionName(secu *spec.Security) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("security_%s_%s", secu.Schema, secu.Name))
}

func (g *GenRoute) GetMiddlewareFunctionName(midd *spec.Middleware) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("middleware_%s", midd.Name))
}
