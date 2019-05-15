package route

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/namecase"
)

func (g *GenRoute) getVarName(name string, typ *spec.Type) string {
	if typ == nil {
		if name == "" {
			return "_"
		}
		return name
	}

	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(typ))), 16)

	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	for typ != nil && (name == "_" || name == "") {
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		name = typ.Name
		typ = typ.Elem
	}

	return g.named.GetSubNamed("").GetName("_"+namecase.ToLowerHumpInitialisms(fmt.Sprintf("%s", name)), addr)
}

func (g *GenRoute) getRouteName(typ *spec.Type) string {
	name := namecase.ToUpperHumpInitialisms(fmt.Sprintf("route_%s", typ.Name))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(typ))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) getOperationFunctionName(oper *spec.Operation) string {
	name := "_" + namecase.ToLowerHumpInitialisms(fmt.Sprintf("operation_%s_%s", oper.Method, oper.Path))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(oper))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) getRequestFunctionName(req *spec.Request) string {
	name := "_" + namecase.ToLowerHumpInitialisms(fmt.Sprintf("request_%s_%s", req.In, req.Name))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(req))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) getSecurityFunctionName(secu *spec.Security) string {
	name := "_" + namecase.ToLowerHumpInitialisms(fmt.Sprintf("security_%s_%s", secu.Schema, secu.Name))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(secu))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) getMiddlewareFunctionName(midd *spec.Middleware) string {
	name := "_" + namecase.ToLowerHumpInitialisms(fmt.Sprintf("middleware_%s_%s", midd.Schema, midd.Name))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(midd))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenRoute) getWrappingFunctionName(wrap *spec.Wrapping) string {
	name := "_" + namecase.ToLowerHumpInitialisms(fmt.Sprintf("wrapping_%s_%s", wrap.Schema, wrap.Name))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(wrap))), 16)
	return g.named.GetName(name, addr)
}
