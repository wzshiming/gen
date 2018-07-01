package parser

import (
	"reflect"
	"strings"

	"github.com/wzshiming/gotype"
)

func GetName(t gotype.Type, tag reflect.StructTag) string {
	name, ok := tag.Lookup("name")
	if !ok {
		name = t.Name()
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
