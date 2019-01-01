package parser

import (
	"fmt"
	"strings"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/gotype"
	"github.com/wzshiming/namecase"
)

// Parser is the parse type generating definitions
type Parser struct {
	imp    *gotype.Importer
	api    *spec.API
	ways   map[string]bool
	consts map[string]gotype.Type
}

func NewParser(imp *gotype.Importer) *Parser {
	return &Parser{
		imp:    imp,
		api:    spec.NewAPI(),
		ways:   map[string]bool{},
		consts: map[string]gotype.Type{},
	}
}

func (g *Parser) API() *spec.API {
	return g.api
}

func (g *Parser) Import(pkgpath string, ways string) error {
	if !strings.HasSuffix(pkgpath, "/...") {
		return g.importOnce(pkgpath)
	}

	if ways != "" {
		for _, way := range strings.Split(ways, ",") {
			g.ways[way] = true
		}
	}

	pkgs := utils.PackageOmitted(pkgpath)
	for _, out := range pkgs {
		err := g.importOnce(out)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Parser) isWay(way string) bool {
	if way == "" {
		return true
	}
	pre := 0
	for i, c := range way {
		if c == ',' {
			if g.ways[way[pre:i]] {
				return true
			}
			pre = i + 1
		}
	}
	return g.ways[way[pre:]]
}

func (g *Parser) importChild(pkgpath, name string) (gotype.Type, bool) {
	t, ok := g.consts[name]
	if ok {
		return t, true
	}

	pkg, err := g.imp.Import(pkgpath)
	if err != nil {
		return nil, false
	}
	return pkg.ChildByName(name)
}

func (g *Parser) importOnce(pkgpath string) error {
	pkg, err := g.imp.Import(pkgpath)
	if err != nil {
		return err
	}
	g.api.Package = pkg.Name()
	numchi := pkg.NumChild()

	for i := 0; i != numchi; i++ {
		v := pkg.Child(i)
		if !IsExported(v.Name()) {
			continue
		}
		switch v.Kind() {
		case gotype.Declaration:

			if len(g.ways) != 0 {
				doc := v.Doc().Text()
				tag := GetTag(doc)
				if !g.isWay(tag.Get("way")) {
					continue
				}
			}

			err = g.AddMiddleware(nil, v)
			if err != nil {
				return err
			}
			err = g.AddSecurity(nil, v)
			if err != nil {
				return err
			}
			err = g.AddOperation("", nil, v)
			if err != nil {
				return err
			}
		default:

			if len(g.ways) != 0 {
				doc := v.Doc().Text()
				tag := GetTag(doc)
				if !g.isWay(tag.Get("way")) {
					continue
				}
			}

			err = g.AddPaths(v)
			if err != nil {
				return err
			}
		case gotype.Interface, gotype.Scope, gotype.Invalid:
			// No action
		}
	}

	return nil
}

func (g *Parser) AddPaths(t gotype.Type) (err error) {
	numm := t.NumMethod()
	if numm == 0 {
		return nil
	}
	doc := strings.TrimSpace(t.Doc().Text())
	if doc == "" {
		return nil
	}
	tag := GetTag(doc)
	path := tag.Get("path")
	if path == "" {
		return nil
	}

	sch, err := g.AddType(t)
	if err != nil {
		return err
	}

	for i := 0; i != numm; i++ {
		v := t.Method(i)
		if !IsExported(v.Name()) {
			continue
		}

		if len(g.ways) != 0 {
			doc := v.Doc().Text()
			tag := GetTag(doc)
			if !g.isWay(tag.Get("way")) {
				continue
			}
		}

		err = g.AddMiddleware(sch, v)
		if err != nil {
			return err
		}
		err = g.AddSecurity(sch, v)
		if err != nil {
			return err
		}
		err = g.AddOperation(path, sch, v)
		if err != nil {
			return err
		}
	}
	return
}

func (g *Parser) AddMiddleware(sch *spec.Type, t gotype.Type) (err error) {
	name := t.Name()
	doc := strings.TrimSpace(t.Doc().Text())
	pkgpath := t.PkgPath()
	if doc == "" {
		return nil
	}
	if t.Kind() == gotype.Declaration {
		t = t.Declaration()
	}
	if t.Kind() != gotype.Func {
		return nil
	}

	tag := GetTag(doc)
	name = GetName(name, tag)
	middleware := tag.Get("middleware")
	if middleware == "" {
		return nil
	}

	path := ""
	route := tag.Get("route")
	if route != "" {
		_, path, _ = GetRoute(route)
	}

	midd := &spec.Middleware{}
	midd.PkgPath = pkgpath
	midd.Schema = middleware
	midd.Type = sch
	midd.Description = doc
	midd.Name = name

	reqs, err := g.AddRequests(path, t)
	midd.Requests = reqs

	resps, err := g.AddResponses(t)
	midd.Responses = resps

	key := name + "." + utils.Hash(name, middleware, doc)

	g.api.Middlewares[key] = midd
	return nil
}

func (g *Parser) AddSecurity(sch *spec.Type, t gotype.Type) (err error) {
	name := t.Name()
	doc := strings.TrimSpace(t.Doc().Text())
	pkgpath := t.PkgPath()
	if doc == "" {
		return nil
	}
	if t.Kind() == gotype.Declaration {
		t = t.Declaration()
	}
	if t.Kind() != gotype.Func {
		return nil
	}

	tag := GetTag(doc)
	name = GetName(name, tag)
	security := tag.Get("security")
	if security == "" {
		return nil
	}

	secu := &spec.Security{}
	secu.PkgPath = pkgpath
	secu.Schema = security
	secu.Type = sch
	secu.Description = doc
	secu.Name = name

	reqs, err := g.AddRequests("", t)
	secu.Requests = reqs

	resps, err := g.AddResponses(t)
	secu.Responses = resps

	key := name + "." + utils.Hash(name, security, doc)

	g.api.Securitys[key] = secu
	return nil
}

func (g *Parser) AddOperation(basePath string, sch *spec.Type, t gotype.Type) (err error) {
	name := t.Name()
	doc := strings.TrimSpace(t.Doc().Text())
	pkgpath := t.PkgPath()

	if doc == "" {
		return nil
	}
	if t.Kind() == gotype.Declaration {
		t = t.Declaration()
	}
	if t.Kind() != gotype.Func {
		return nil
	}
	tag := GetTag(doc)
	name = GetName(name, tag)
	route := tag.Get("route")
	if route == "" {
		return nil
	}
	method, pat, ok := GetRoute(route)
	if !ok {
		return nil
	}

	oper := &spec.Operation{}
	if basePath != "" {
		oper.Tags = append(oper.Tags, namecase.ToCamel(basePath))
		pat = Join(basePath, pat)
	}
	oper.PkgPath = pkgpath
	oper.Method = method
	oper.Path = pat
	oper.Description = doc
	oper.Type = sch
	oper.Name = name

	reqs, err := g.AddRequests(pat, t)
	oper.Requests = reqs

	resps, err := g.AddResponses(t)
	oper.Responses = resps

	g.api.Operations = append(g.api.Operations, oper)
	return nil
}

func (g *Parser) AddResponses(t gotype.Type) (resps []*spec.Response, err error) {
	numout := t.NumOut()
	for i := 0; i != numout; i++ {
		v := t.Out(i)
		resp, err := g.AddResponse(v)
		if err != nil {
			return nil, err
		}
		resps = append(resps, resp)
	}
	return resps, nil
}

func (g *Parser) AddResponse(t gotype.Type) (resp *spec.Response, err error) {

	if t.Kind() != gotype.Declaration {
		return nil, fmt.Errorf("Gen.AddResponse: unsupported type: %s is %s kind\n", t.String(), t.Kind().String())
	}
	name := t.Name()
	doc := strings.TrimSpace(t.Comment().Text())
	tag := GetTag(doc)
	name = GetName(name, tag)
	code := tag.Get("code")
	in := tag.Get("in")
	content := tag.Get("content")
	t = t.Declaration()

	kind := t.Kind()
	if in == "" {
		in = "body"
	}

	if in == "body" {
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
	}

	sch, err := g.AddType(t)
	if err != nil {
		return nil, err
	}

	key := name + "." + utils.Hash(in, sch.Name, sch.Ref, doc)

	if g.api.Responses[key] != nil {
		return &spec.Response{
			Ref: key,
		}, nil
	}

	resp = &spec.Response{}
	resp.Name = name
	resp.In = in
	resp.Code = code
	resp.Content = content
	resp.Type = sch
	resp.Description = doc

	g.api.Responses[key] = resp
	return &spec.Response{
		Ref: key,
	}, nil
}

func (g *Parser) AddRequests(path string, t gotype.Type) (reqs []*spec.Request, err error) {
	numin := t.NumIn()
	for i := 0; i != numin; i++ {
		v := t.In(i)
		req, err := g.AddRequest(path, v)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (g *Parser) AddRequest(path string, t gotype.Type) (par *spec.Request, err error) {

	if t.Kind() != gotype.Declaration {
		return nil, fmt.Errorf("Gen.AddRequest: unsupported type: %s is %s kind\n", t.String(), t.Kind().String())
	}
	name := t.Name()
	doc := strings.TrimSpace(t.Comment().Text())
	tag := GetTag(doc)
	name = GetName(name, tag)
	in := tag.Get("in")
	t = t.Declaration()

	if in == "" {
		tt := t
		ps := 0
		for tt.Kind() == gotype.Ptr {
			tt = tt.Elem()
			ps++
		}

		tname := tt.Name()
		tpkgpath := tt.PkgPath()

		switch tpkgpath {
		case "net/http":
			switch tname {
			case "Request":
				if ps == 1 {
					return &spec.Request{
						In:   "none",
						Name: fmt.Sprintf("*%s.%s", tpkgpath, tname),
					}, nil
				}
			case "ResponseWriter":
				return &spec.Request{
					In:   "none",
					Name: fmt.Sprintf("%s.%s", tpkgpath, tname),
				}, nil
			}
		}

		switch tt.Kind() {
		case gotype.Array, gotype.Slice, gotype.Map, gotype.Struct:
			in = "body"
		default:
			if path == "" || strings.Index(path, "{"+name+"}") == -1 {
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

	sch, err := g.AddType(t)
	if err != nil {
		return nil, err
	}

	key := name + "." + utils.Hash(in, sch.Name, sch.Ref, doc)

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
		Ref: key,
	}, nil
}

func (g *Parser) AddType(t gotype.Type) (sch *spec.Type, err error) {
	name := t.Name()
	pkgpath := t.PkgPath()
	doc := strings.TrimSpace(t.Doc().Text())
	tag := GetTag(doc)
	name = GetName(name, tag)
	kind := t.Kind()

	key := name + "." + utils.Hash(name, pkgpath, t.String(), doc)
	if g.api.Types[key] != nil {
		return &spec.Type{
			Ref: key,
		}, nil
	}

	sch = &spec.Type{}
	sch.PkgPath = pkgpath
	sch.Name = name
	sch.Description = doc

	if time, ok := g.importChild("time", "Time"); ok && gotype.Equal(time, t) {
		sch.Description = "This is the time string in RFC3339 format"
		sch.Kind = spec.Time
		return sch, nil
	}

	switch kind {
	case gotype.Struct:
		// Field
		{
			num := t.NumField()
			for i := 0; i != num; i++ {
				v := t.Field(i)
				name := v.Name()
				if !IsExported(name) {
					continue
				}
				tag := v.Tag()
				val, err := g.AddType(v.Elem())
				if err != nil {
					return nil, err
				}

				desc := v.Doc().Text() + "\n" + v.Comment().Text()
				desc = strings.TrimSpace(desc)
				field := &spec.Field{
					Name:        name,
					Type:        val,
					Tag:         tag,
					Anonymous:   v.IsAnonymous(),
					Description: desc,
				}

				sch.Fields = append(sch.Fields, field)
			}
		}

	case gotype.Error, gotype.String, gotype.Bool, gotype.Float32, gotype.Float64,
		gotype.Int8, gotype.Int16, gotype.Int32, gotype.Int64, gotype.Int,
		gotype.Uint8, gotype.Uint16, gotype.Uint32, gotype.Uint64, gotype.Uint,
		gotype.Byte, gotype.Rune:

		if name != "_" && name != strings.ToLower(kind.String()) {
			scope, err := g.imp.Import(t.PkgPath())
			if err != nil {
				return nil, err
			}

			numchi := scope.NumChild()
			for i := 0; i != numchi; i++ {
				v := scope.Child(i)
				if v.Kind() != gotype.Declaration {
					continue
				}
				v = v.Declaration()
				if typname := v.Name(); name == typname {
					if value := v.Value(); value != "" {
						desc := v.Doc().Text() + "\n" + v.Comment().Text()
						desc = strings.TrimSpace(desc)
						sch.Enum = append(sch.Enum, &spec.Enum{
							Value:       value,
							Description: desc,
						})
					}
				}
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
		sch.Key = schk
		sch.Elem = schv
	case gotype.Slice:
		schv, err := g.AddType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch.Elem = schv
	case gotype.Array:
		schv, err := g.AddType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch.Elem = schv
		sch.Len = t.Len()
	case gotype.Ptr:
		schv, err := g.AddType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch.Elem = schv
	default:
		return nil, fmt.Errorf("Gen.AddType: unsupported type: %s is %s kind\n", t.String(), t.Kind().String())
	}

	sch.Kind = kindMapping[kind]

	if name != "" && name != strings.ToLower(kind.String()) {
		g.api.Types[key] = sch
		return &spec.Type{
			Ref: key,
		}, nil
	}

	return sch, nil
}
