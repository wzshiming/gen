package spec

import (
	"reflect"
)

type API struct {
	Package     string
	Operations  []*Operation
	Middlewares map[string]*Middleware
	Securitys   map[string]*Security
	Requests    map[string]*Request
	Responses   map[string]*Response
	Types       map[string]*Type
}

func NewAPI() *API {
	return &API{
		Middlewares: map[string]*Middleware{},
		Securitys:   map[string]*Security{},
		Requests:    map[string]*Request{},
		Responses:   map[string]*Response{},
		Types:       map[string]*Type{},
	}
}

type Middleware struct {
	Ref         string
	PkgPath     string
	Schema      string
	Name        string
	Type        *Type
	Requests    []*Request
	Responses   []*Response
	Description string
}

type Security struct {
	//	Ref         string
	PkgPath     string
	Schema      string
	Name        string
	Type        *Type
	Requests    []*Request
	Responses   []*Response
	Description string
}

type Operation struct {
	PkgPath     string
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
	Ident       string
	Name        string
	In          string
	Content     string
	Type        *Type
	Description string
}

type Response struct {
	Ref         string
	Ident       string
	Name        string
	In          string
	Code        string
	Content     string
	Type        *Type
	Description string
}

type Type struct {
	Ref         string
	Ident       string
	PkgPath     string
	Name        string
	Kind        Kind
	Key         *Type
	Elem        *Type
	Fields      []*Field
	Len         int
	Enum        []*Enum
	Description string

	IsRoot            bool
	IsJSONUnmarshaler bool
	IsJSONMarshaler   bool
	IsTextUnmarshaler bool
	IsTextMarshaler   bool
	IsReader          bool
	IsImage           bool
}

type Field struct {
	Name        string
	Type        *Type
	Tag         reflect.StructTag
	Anonymous   bool
	Description string
}

type Enum struct {
	Value       string
	Description string
}
