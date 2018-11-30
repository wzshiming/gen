package route

import (
	"fmt"

	"github.com/wzshiming/namecase"
)

func GetGlobalVarName(name string) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("global_%s", name))
}

func GetRouteName(name string) string {
	return namecase.ToPascal(fmt.Sprintf("route_%s", name))
}

func GetOperationFunctionName(name string) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("operation_%s", name))
}

func GetRequestFunctionName(name, in string) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("request_%s_%s", in, name))
}

func GetSecurityFunctionName(name string) string {
	return "_" + namecase.ToCamel(fmt.Sprintf("security_%s", name))
}
