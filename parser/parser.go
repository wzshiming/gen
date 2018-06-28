package parser

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/gotype"
)

// Parser is the parse type generating definitions
type Parser struct {
	imp *gotype.Importer
	api *spec.API
}

func NewParser() *Parser {
	return &Parser{
		imp: gotype.NewImporter(gotype.WithCommentLocator()),
		api: spec.NewAPI(),
	}
}

func (g *Parser) API() *spec.API {
	return g.api
}

func (g *Parser) Import(pkgpath string) error {
	pkg, err := g.imp.Import(pkgpath)
	if err != nil {
		return err
	}
	g.api.Package = pkg.Name()
	numchi := pkg.NumChild()

	for i := 0; i != numchi; i++ {
		v := pkg.Child(i)
		switch v.Kind() {
		case gotype.Func:
			err := g.AddOperation("", nil, v)
			if err != nil {
				return err
			}
		default:
			err := g.AddPaths(v)
			if err != nil {
				return err
			}
		case gotype.Interface, gotype.Scope, gotype.Invalid, gotype.Var:
			// No action
		}
	}

	return nil
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

func (g *Parser) AddPaths(t gotype.Type) (err error) {
	numm := t.NumMethods()
	if numm == 0 {
		return nil
	}
	tag := GetTag(t.Doc().Text())
	path := tag.Get("path")
	if path == "" {
		return nil
	}

	sch, err := g.AddType(t)
	if err != nil {
		return err
	}
	for i := 0; i != numm; i++ {
		v := t.Methods(i)
		err := g.AddOperation(path, sch, v)
		if err != nil {
			return err
		}
	}
	return
}

func (g *Parser) AddOperation(basePath string, sch *spec.Type, t gotype.Type) (err error) {
	if t.Kind() != gotype.Func {
		return fmt.Errorf("Gen: unsupported type: %s", t.Kind().String())
	}

	doc := t.Doc().Text()
	tag := GetTag(doc)
	route := tag.Get("route")
	if route == "" {
		return nil
	}
	rs := strings.SplitN(route, " ", 2)
	if len(rs) != 2 {
		return nil
	}

	pat := strings.TrimSpace(rs[1])

	method := strings.ToLower(strings.TrimSpace(rs[0]))

	oper := &spec.Operation{}
	if basePath != "" {
		oper.Tags = append(oper.Tags, basePath)
		pat = path.Join(basePath, pat)
	}
	oper.Method = method
	oper.Path = pat
	oper.Description = doc
	oper.Type = sch

	oper.Name = t.Name()
	{
		numin := t.NumIn()
		for i := 0; i != numin; i++ {
			v := t.In(i)
			par, err := g.AddRequest(pat, v)
			if err != nil {
				return err
			}
			if par != nil {
				oper.Requests = append(oper.Requests, par)
			}
		}
	}

	{
		numout := t.NumOut()
		for i := 0; i != numout; i++ {
			v := t.Out(i)
			resp, err := g.AddResponse(v)
			if err != nil {
				return err
			}
			oper.Responses = append(oper.Responses, resp)
		}
	}
	g.api.Operations = append(g.api.Operations, oper)
	return nil
}

func (g *Parser) AddResponse(t gotype.Type) (resp *spec.Response, err error) {
	doc := t.Comment().Text()
	tag := GetTag(doc)
	name := t.Name()
	code := tag.Get("code")
	content := tag.Get("content")
	kind := t.Elem().Kind()
	if code == "" {
		if kind != gotype.Error {
			code = "200"
		} else {
			code = "400"
		}
	}

	if content == "" {
		if kind != gotype.Error {
			content = "json"
		} else {
			content = "error"
		}
	}

	sch, err := g.AddType(t.Elem())
	if err != nil {
		return nil, err
	}

	key := name + "." + utils.Hash(name, code, content, sch.Name, doc)

	if g.api.Responses[key] != nil {
		return &spec.Response{
			Ref: key,
		}, nil
	}

	resp = &spec.Response{}
	resp.Name = name
	resp.Code = code
	resp.Content = content
	resp.Type = sch
	resp.Description = doc

	g.api.Responses[key] = resp
	return &spec.Response{
		Ref:  key,
		Name: sch.Name,
	}, nil
}

func (g *Parser) AddRequest(path string, t gotype.Type) (par *spec.Request, err error) {
	rawname := t.Name()
	doc := t.Comment().Text()
	tag := GetTag(doc)

	in := tag.Get("in")
	if in == "" {
		t := t
		for t.Elem().Kind() == gotype.Ptr {
			t = t.Elem()
		}
		switch t.Elem().Kind() {
		case gotype.Array, gotype.Slice, gotype.Map, gotype.Struct:
			in = "body"
		default:
			if strings.Index(path, "{"+rawname+"}") == -1 {
				in = "query"
			} else {
				in = "path"
			}
		}
	}

	content := tag.Get("content")
	if content == "" && in == "body" {
		content = "json"
	}

	name, ok := tag.Lookup("name")
	if !ok {
		name = rawname
	}

	sch, err := g.AddType(t.Elem())
	if err != nil {
		return nil, err
	}

	key := name + "." + utils.Hash(name, in, sch.Name, doc)

	if g.api.Requests[key] != nil {
		return &spec.Request{
			Ref: key,
		}, nil
	}
	par = &spec.Request{}
	par.In = in
	par.Name = name
	par.Content = content
	par.Description = doc
	par.Type = sch

	g.api.Requests[key] = par
	return &spec.Request{
		Ref:  key,
		Name: sch.Name,
	}, nil
}

func (g *Parser) AddType(t gotype.Type) (sch *spec.Type, err error) {
	name := t.Name()
	pkgpath := t.PkgPath()
	doc := t.Doc().Text()
	kind := t.Kind()

	key := name + "." + utils.Hash(name, pkgpath, kind.String(), doc)
	if g.api.Types[key] != nil {
		return &spec.Type{
			Ref: key,
		}, nil
	}

	if t.IsGoroot() && pkgpath == "time" && name == "Time" {
		sch := &spec.Type{
			Name: name,
			Kind: spec.Time,
		}
		return sch, nil
	}

	switch kind {
	case gotype.Struct:

		sch = &spec.Type{}

		// Field
		{
			num := t.NumField()
			for i := 0; i != num; i++ {
				v := t.Field(i)
				name := v.Name()
				tag := v.Tag()
				val, err := g.AddType(v.Elem())
				if err != nil {
					return nil, err
				}
				field := &spec.Field{
					Name:        name,
					Type:        val,
					Tag:         tag,
					Anonymous:   v.IsAnonymous(),
					Description: v.Doc().Text() + v.Comment().Text(),
				}

				sch.Fields = append(sch.Fields, field)
			}
		}

	case gotype.Error, gotype.String, gotype.Bool, gotype.Float32, gotype.Float64,
		gotype.Int8, gotype.Int16, gotype.Int32, gotype.Int64, gotype.Int,
		gotype.Uint8, gotype.Uint16, gotype.Uint32, gotype.Uint64, gotype.Uint,
		gotype.Byte, gotype.Rune:
		sch = &spec.Type{}

		scope, err := g.imp.Import(t.PkgPath())
		if err != nil {
			return nil, err
		}

		typname := t.Name()
		numchi := scope.NumChild()
		for i := 0; i != numchi; i++ {
			v := scope.Child(i)
			if v.Kind() != gotype.Var {
				continue
			}
			if v.Elem().Name() == typname {
				name := v.Name()
				if name == "_" {
					continue
				}
				value := v.Value()
				sch.Enum = append(sch.Enum, value)
			}
		}

	case gotype.Map:
		schk, err := g.AddType(t.Key())
		if err != nil {
			return nil, err
		}
		schv, err := g.AddType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch = &spec.Type{
			Key:  schk,
			Elem: schv,
		}
	case gotype.Slice:
		sch, err = g.AddType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch = &spec.Type{
			Elem: sch,
		}
	case gotype.Array:
		sch, err = g.AddType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch = &spec.Type{
			Elem: sch,
			Len:  t.Len(),
		}
	case gotype.Ptr:
		sch, err = g.AddType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch = &spec.Type{
			Elem: sch,
		}
	default:
		return nil, fmt.Errorf("Gen: unsupported type: %s", t.Kind().String())
	}

	sch.Name = name
	sch.Kind = kindMapping[kind]
	sch.Description = doc

	//	tag := GetTag(doc)
	//	if typ := tag.Get("type"); typ != "" {
	//		sch.Kind = typ
	//	}
	if name != "" && name != strings.ToLower(kind.String()) {
		g.api.Types[key] = sch
		return &spec.Type{
			Ref:  key,
			Name: sch.Name,
		}, nil
	}

	return sch, nil
}
