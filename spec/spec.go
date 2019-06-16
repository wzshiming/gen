package spec

import (
	"reflect"
)

type API struct {
	Imports     []string
	Package     string
	Operations  []*Operation
	Middlewares map[string]*Middleware
	Wrappings   map[string]*Wrapping
	Securitys   map[string]*Security
	Requests    map[string]*Request
	Responses   map[string]*Response
	Types       map[string]*Type
}

func NewAPI() *API {
	return &API{
		Middlewares: map[string]*Middleware{},
		Wrappings:   map[string]*Wrapping{},
		Securitys:   map[string]*Security{},
		Requests:    map[string]*Request{},
		Responses:   map[string]*Response{},
		Types:       map[string]*Type{},
	}
}

type Middleware struct {
	Ref            string
	PkgPath        string
	Schema         string
	Name           string
	Type           *Type
	Requests       []*Request
	Responses      []*Response
	Description    string
	DescriptionTag reflect.StructTag
}

type Wrapping struct {
	PkgPath        string
	Schema         string
	Name           string
	Type           *Type
	Requests       []*Request
	Responses      []*Response
	Description    string
	DescriptionTag reflect.StructTag
}

type Security struct {
	//	Ref         string
	PkgPath        string
	Schema         string
	Name           string
	Type           *Type
	Requests       []*Request
	Responses      []*Response
	Description    string
	DescriptionTag reflect.StructTag
}

type Operation struct {
	PkgPath   string
	Method    string
	BasePath  string
	Path      string
	Tags      []string
	Name      string
	Type      *Type
	Chain     []string
	Requests  []*Request
	Responses []*Response
	Securitys []*Security

	Deprecated     bool
	Summary        string
	Description    string
	DescriptionTag reflect.StructTag
}

type Request struct {
	Ref            string
	Ident          string
	Name           string
	In             string
	Content        string
	Type           *Type
	Description    string
	DescriptionTag reflect.StructTag
}

type Response struct {
	Ref            string
	Ident          string
	Name           string
	In             string
	Code           string
	Content        string
	Type           *Type
	Description    string
	DescriptionTag reflect.StructTag
}

type Type struct {
	Ref            string
	Ident          string
	PkgPath        string
	Name           string
	Kind           Kind
	Key            *Type
	Elem           *Type
	Fields         []*Field
	Len            int
	Enum           []*Enum
	Description    string
	DescriptionTag reflect.StructTag
	Attr           Attr
}

type Field struct {
	Name           string
	Type           *Type
	Tag            reflect.StructTag
	Anonymous      bool
	Description    string
	DescriptionTag reflect.StructTag
}

type Enum struct {
	Name           string
	Value          string
	Description    string
	DescriptionTag reflect.StructTag
}

type Attr uint64

const (
	AttrRoot Attr = 1 << (63 - iota)
	AttrJSONUnmarshaler
	AttrJSONMarshaler
	AttrTextUnmarshaler
	AttrTextMarshaler
	AttrReader
	AttrImage
)

func (a Attr) HasOne(b Attr) bool {
	return a&b != 0
}

func (a Attr) Has(b Attr) bool {
	return a&b == b
}

func (a *Attr) Add(b Attr) {
	*a |= b
}
