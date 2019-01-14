package client

import (
	"fmt"

	"github.com/wzshiming/namecase"
)

func (g *GenClient) GetVarName(name string) string {
	return "_" + namecase.ToLowerHumpInitialisms(fmt.Sprintf("var_%s", name))
}
