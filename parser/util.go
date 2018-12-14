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

// GetTag [#[^#]+#]...
func GetTag(text string) reflect.StructTag {
	ss := []string{}
	prev := 0
	for i, v := range text {
		if v != '#' {
			continue
		}
		if prev == 0 {
			prev = i
		} else {
			ss = append(ss, text[prev+1:i])
			prev = 0
		}
	}
	return reflect.StructTag(strings.Join(ss, " "))
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
