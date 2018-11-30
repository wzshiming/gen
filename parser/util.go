package parser

import (
	"go/ast"
	"reflect"
	"strings"
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
