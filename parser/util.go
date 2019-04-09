package parser

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/gotype"
)

func GetName(t string, tag reflect.StructTag) string {
	name, ok := tag.Lookup("name")
	if !ok {
		name = t
	}
	return name
}

func IsExported(name string) bool {
	return ast.IsExported(name)
}

func GetTypeHash(typ gotype.Type) string {
	tp := 0
	for typ.Kind() == gotype.Ptr {
		tp++
		typ = typ.Elem()
	}

	pkgpath := typ.PkgPath()
	name := typ.Name()
	return strings.Repeat("_", tp) + name + "." + utils.Hash(name, pkgpath)
}

// GetRoute "!GET /path"
func GetRoute(route string) (deprecated bool, method, path string, ok bool) {

	route = strings.TrimSpace(route)
	if route == "" {
		return false, "", "", false
	}
	if route[0] == '!' {
		deprecated = true
		route = route[1:]
	}

	route = strings.TrimSpace(route)
	rs := strings.SplitN(route, " ", 2)
	if len(rs) != 2 {
		return false, "", "", false
	}
	path = strings.TrimSpace(rs[1])
	if path == "" {
		return false, "", "", false
	}
	method = strings.TrimSpace(strings.ToUpper(rs[0]))
	if method == "" {
		return false, "", "", false
	}
	return deprecated, method, path, true
}
