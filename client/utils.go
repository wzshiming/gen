package client

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/namecase"
)

var errResponse = &spec.Response{
	Name: "err",
	Type: &spec.Type{
		Kind: spec.Error,
	},
}

func (g *GenClient) getVarName(name string, typ *spec.Type) string {
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
	if typ.Kind == spec.Error {
		return "err"
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

func (g *GenClient) getTypeName(typ *spec.Type) string {
	name := typ.Name
	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(typ))), 16)
	return g.named.GetName(name, addr)
}

func (g *GenClient) getFuncName(oper *spec.Operation) string {
	name := namecase.ToUpperHumpInitialisms(strings.Join(append(oper.Chain, oper.Name), "_"))

	named := g.named
	if oper.Type != nil {
		addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(oper.Type))), 16)
		named = named.GetSubNamed(addr)
	}

	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(oper))), 16)
	return named.GetName(name, addr)
}

func (g *GenClient) getSecurityName(secu *spec.Security) string {
	name := secu.Name

	named := g.named
	if secu.Type != nil {
		addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(secu.Type))), 16)
		named = named.GetSubNamed(addr)
	}

	addr := strconv.FormatUint(uint64(uintptr(unsafe.Pointer(secu))), 16)
	return named.GetName(name, addr)
}

func (g *GenClient) getEnumName(name, value string) string {
	return g.named.GetName(name, value)
}
