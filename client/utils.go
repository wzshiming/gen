package client

import (
	"fmt"

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
	return "_" + namecase.ToLowerHumpInitialisms(fmt.Sprintf("var_%s", name))
}
