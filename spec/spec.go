package spec

import (
	"reflect"
)

type API struct {
	Package    string
	Operations []*Operation
	Securitys  map[string]*Security
	Requests   map[string]*Request
	Responses  map[string]*Response
	Types      map[string]*Type
}

func NewAPI() *API {
	return &API{
		Securitys: map[string]*Security{},
		Requests:  map[string]*Request{},
		Responses: map[string]*Response{},
		Types:     map[string]*Type{},
	}
}

type Security struct {
	//	Ref         string
	Schema      string
	Name        string
	Type        *Type
	Requests    []*Request
	Responses   []*Response
	Description string
}

type Operation struct {
	Method      string
	Path        string
	Tags        []string
	Name        string
	Type        *Type
	Requests    []*Request
	Responses   []*Response
	Securitys   []*Security
	Description string
}

type Request struct {
	Ref         string
	Name        string
	In          string
	Content     string
	Type        *Type
	Description string
}

type Response struct {
	Ref         string
	Name        string
	Code        string
	Content     string
	Type        *Type
	Description string
}

type Type struct {
	Ref         string
	Name        string
	Kind        Kind
	Key         *Type
	Elem        *Type
	Fields      []*Field
	Len         int
	Enum        []string
	Description string
}

type Field struct {
	Name        string
	Type        *Type
	Tag         reflect.StructTag
	Anonymous   bool
	Description string
}
