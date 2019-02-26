package client

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/namecase"
)

func (g *GenClient) getVarName(name string, typ *spec.Type) string {
	if typ == nil {
		if name == "" {
			return "_"
		}
		return name
	}
	if typ.Kind == spec.Error {
		return "err"
	}
	if name == "" {
		name = typ.Name
	}
	return "_" + namecase.ToLowerHumpInitialisms(name)
}

func (g *GenClient) getTypeName(typ *spec.Type) string {
	name := namecase.ToUpperHumpInitialisms(typ.Name)
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(typ))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenClient) getFuncName(oper *spec.Operation) string {
	name := namecase.ToUpperHumpInitialisms(fmt.Sprintf("%s_%s", oper.Method, oper.Path))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(oper))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenClient) getSecurityName(secu *spec.Security) string {
	name := namecase.ToLowerHumpInitialisms(fmt.Sprintf("security_%s_%s", secu.Schema, secu.Name))
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(secu))), 16)
	return g.named.GetName(name, addr)
}
