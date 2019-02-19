package parser

import (
	"fmt"
	"path"
	"strings"

	"github.com/wzshiming/gen/named"
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/gotype"
)

// Parser is the parse type generating definitions
type Parser struct {
	imp    *gotype.Importer
	api    *spec.API
	ways   map[string]bool
	consts map[string]gotype.Type

	namedReq  *named.Named
	namedResp *named.Named
	namedMidd *named.Named
	namedSecu *named.Named
	namedTyp  *named.Named
}

func NewParser(imp *gotype.Importer) *Parser {
	return &Parser{
		imp:       imp,
		api:       spec.NewAPI(),
		ways:      map[string]bool{},
		consts:    map[string]gotype.Type{},
		namedReq:  named.NewNamed("."),
		namedResp: named.NewNamed("."),
		namedMidd: named.NewNamed("."),
		namedSecu: named.NewNamed("."),
		namedTyp:  named.NewNamed("."),
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
	oname := t.Name()
	doc := strings.TrimSpace(t.Doc().Text())
	pkgpath := t.PkgPath()
	if doc == "" {
		return nil
	}
	t = t.Declaration()
	if t.Kind() != gotype.Func {
		return nil
	}

	tag := GetTag(doc)
	name := GetName(oname, tag)
	middleware := tag.Get("middleware")
	if middleware == "" {
		return nil
	}

	path := ""
	route := tag.Get("route")
	if route != "" {
		_, _, path, _ = GetRoute(route)
	}

	hash := utils.Hash(oname, pkgpath, middleware, doc)
	key := g.namedMidd.GetName(name, hash)

	midd := &spec.Middleware{}
	midd.Name = name
	midd.PkgPath = pkgpath
	midd.Schema = middleware
	midd.Description = doc
	midd.Type = sch

	reqs, err := g.AddRequests(path, t)
	midd.Requests = reqs

	resps, err := g.AddResponses(t)
	midd.Responses = resps

	g.api.Middlewares[key] = midd
	return nil
}

func (g *Parser) AddSecurity(sch *spec.Type, t gotype.Type) (err error) {
	oname := t.Name()
	doc := strings.TrimSpace(t.Doc().Text())
	pkgpath := t.PkgPath()
	if doc == "" {
		return nil
	}
	t = t.Declaration()
	if t.Kind() != gotype.Func {
		return nil
	}

	tag := GetTag(doc)
	name := GetName(oname, tag)
	security := tag.Get("security")
	if security == "" {
		return nil
	}

	hash := utils.Hash(oname, pkgpath, security, doc)
	key := g.namedSecu.GetName(name, hash)

	secu := &spec.Security{}
	secu.Name = name
	secu.PkgPath = pkgpath
	secu.Schema = security
	secu.Description = doc
	secu.Type = sch

	reqs, err := g.AddRequests("", t)
	secu.Requests = reqs

	resps, err := g.AddResponses(t)
	secu.Responses = resps

	g.api.Securitys[key] = secu
	return nil
}

func (g *Parser) AddOperation(basePath string, sch *spec.Type, t gotype.Type) (err error) {
	oname := t.Name()
	doc := strings.TrimSpace(t.Doc().Text())
	pkgpath := t.PkgPath()

	if doc == "" {
		return nil
	}
	t = t.Declaration()
	if t.Kind() != gotype.Func {
		return nil
	}
	tag := GetTag(doc)
	name := GetName(oname, tag)
	route := tag.Get("route")
	if route == "" {
		return nil
	}
	deprecated, method, pat, ok := GetRoute(route)
	if !ok {
		return nil
	}

	if basePath != "" {
		basePath = path.Clean(basePath)
		pat = path.Join(basePath, pat)
	}

	oper := &spec.Operation{}
	oper.PkgPath = pkgpath
	oper.Method = method
	oper.BasePath = basePath
	oper.Path = pat
	oper.Description = doc
	oper.Summary = strings.SplitN(oper.Description, "\n", 2)[0]
	oper.Deprecated = deprecated
	oper.Type = sch
	oper.Name = name

	if sch != nil {
		sch := sch
		if sch.Ref != "" {
			sch = g.api.Types[sch.Ref]
		}
		oper.Tags = append(oper.Tags, sch.Name)
	}

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

	oname := t.Name()
	doc := strings.TrimSpace(t.Comment().Text())
	tag := GetTag(doc)
	name := GetName(oname, tag)
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

	si := sch.Ref
	if si == "" {
		si = sch.Ident
	}
	hash := utils.Hash(oname, in, code, content, doc, si)
	key := g.namedResp.GetName(name, hash)

	if g.api.Responses[key] != nil {
		return &spec.Response{
			Ref: key,
		}, nil
	}

	resp = &spec.Response{}
	resp.Ident = key
	resp.In = in
	resp.Name = name
	resp.Code = code
	resp.Content = content
	resp.Description = doc
	resp.Type = sch

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

	oname := t.Name()
	doc := strings.TrimSpace(t.Comment().Text())
	tag := GetTag(doc)
	name := GetName(oname, tag)
	in := tag.Get("in")
	t = t.Declaration()
	switch t.Kind() {
	case gotype.Ptr:
		if req, ok := g.importChild("net/http", "Request"); ok && gotype.Equal(t.Elem(), req) {
			return &spec.Request{
				In:    "none",
				Ident: "*net/http.Request",
				Name:  "r",
			}, nil
		}

	case gotype.Interface:
		if resp, ok := g.importChild("net/http", "ResponseWriter"); ok && gotype.Implements(resp, t) {
			return &spec.Request{
				In:    "none",
				Ident: "net/http.ResponseWriter",
				Name:  "w",
			}, nil
		}
	}

	sch, err := g.AddType(t)
	if err != nil {
		return nil, err
	}

	content := tag.Get("content")

	if in == "" {
		typ := sch
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}

		if typ.Attr.Has(spec.AttrImage) {
			in = "body"
			if content == "" {
				content = "image"
			}
		} else if typ.Attr.Has(spec.AttrReader) {
			in = "body"
			if content == "" {
				content = "file"
			}
		} else if typ.Attr.Has(spec.AttrTextUnmarshaler) {
			if path != "" && strings.Index(path, "{"+name+"}") != -1 {
				in = "path"
			} else {
				in = "query"
			}
		} else {
			if typ.Kind == spec.Ptr {
				typ = typ.Elem
				if typ.Ref != "" {
					typ = g.api.Types[typ.Ref]
				}
			}
			switch typ.Kind {
			case spec.Array, spec.Slice, spec.Map, spec.Struct:
				in = "body"
			case spec.Interface:
				in = "middleware"
			default:
				if path != "" && strings.Index(path, "{"+name+"}") != -1 {
					in = "path"
				} else {
					in = "query"
				}
			}
		}
	}

	if content == "" && in == "body" {
		content = "json"
	}

	si := sch.Ref
	if si == "" {
		si = sch.Ident
	}
	hash := utils.Hash(oname, in, content, doc, si)
	key := g.namedReq.GetName(name, hash)

	if g.api.Requests[key] != nil {
		return &spec.Request{
			Ref: key,
		}, nil
	}

	par = &spec.Request{}
	par.Ident = key
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
	oname := t.Name()
	pkgpath := t.PkgPath()
	doc := strings.TrimSpace(t.Doc().Text())
	tag := GetTag(doc)
	name := GetName(oname, tag)
	kind := t.Kind()
	isRoot := t.IsGoroot()
	isBuiltin := name == strings.ToLower(kind.String())
	hash := ""
	if !isBuiltin {
		hash = utils.Hash(t.String(), oname, pkgpath, doc)
	}
	key := g.namedTyp.GetName(name, hash)

	if g.api.Types[key] != nil {
		return &spec.Type{
			Ref: key,
		}, nil
	}

	sch = &spec.Type{}
	if isRoot {
		sch.Attr.Add(spec.AttrRoot)
	}
	sch.Ident = key
	sch.Name = name
	sch.PkgPath = pkgpath
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
	case gotype.Interface:
		// No action
	default:
		return nil, fmt.Errorf("Gen.AddType: unsupported type: %s is %s kind\n", t.String(), t.Kind().String())
	}

	sch.Kind = kindMapping[kind]

	if text, ok := g.importChild("encoding", "TextUnmarshaler"); ok && gotype.Implements(t, text) {
		sch.Attr.Add(spec.AttrTextUnmarshaler)
	}
	if text, ok := g.importChild("encoding", "TextMarshaler"); ok && gotype.Implements(t, text) {
		sch.Attr.Add(spec.AttrTextMarshaler)
	}
	if text, ok := g.importChild("encoding/json", "Unmarshaler"); ok && gotype.Implements(t, text) {
		sch.Attr.Add(spec.AttrJSONUnmarshaler)
	}
	if text, ok := g.importChild("encoding/json", "Marshaler"); ok && gotype.Implements(t, text) {
		sch.Attr.Add(spec.AttrJSONMarshaler)
	}

	if read, ok := g.importChild("io", "Reader"); ok && gotype.Implements(t, read) {
		sch.Attr.Add(spec.AttrReader)
	}

	if read, ok := g.importChild("image", "Image"); ok && gotype.Implements(t, read) {
		sch.Attr.Add(spec.AttrImage)
	}

	if name != "" && !isBuiltin {
		g.api.Types[key] = sch
		return &spec.Type{
			Ref: key,
		}, nil
	}

	return sch, nil
}
