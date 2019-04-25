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
	namedWrap *named.Named
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
		namedWrap: named.NewNamed("."),
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
	g.api.Imports = append(g.api.Imports, pkgpath)
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
				_, tag := utils.GetTag(doc)
				if !g.isWay(tag.Get("way")) {
					continue
				}
			}

			err = g.addMiddleware(nil, v)
			if err != nil {
				return err
			}
			err = g.addSecurity(nil, v)
			if err != nil {
				return err
			}
			err = g.addOperation("", nil, v, nil)
			if err != nil {
				return err
			}
		default:

			if len(g.ways) != 0 {
				doc := v.Doc().Text()
				_, tag := utils.GetTag(doc)
				if !g.isWay(tag.Get("way")) {
					continue
				}
			}

			err = g.addPaths(v)
			if err != nil {
				return err
			}
		case gotype.Scope, gotype.Invalid:
			// No action
		}
	}

	return nil
}

func (g *Parser) addPaths(t gotype.Type) (err error) {
	doc := strings.TrimSpace(t.Doc().Text())
	if doc == "" {
		return nil
	}
	_, tag := utils.GetTag(doc)
	path := tag.Get("path")
	if path == "" {
		return nil
	}

	sch, err := g.addType(t)
	if err != nil {
		return err
	}

	filter := map[string]bool{}
	return g.addMethods(path, sch, t, nil, filter)
}

func (g *Parser) addMethods(basePath string, sch *spec.Type, t gotype.Type, chain []string, filter map[string]bool) (err error) {
	numm := t.NumMethod()
	for i := 0; i != numm; i++ {
		v := t.Method(i)
		name := v.Name()
		if !IsExported(name) || filter[name] {
			continue
		}
		filter[name] = true

		if len(g.ways) != 0 {
			doc := v.Doc().Text()
			_, tag := utils.GetTag(doc)
			if !g.isWay(tag.Get("way")) {
				continue
			}
		}

		err = g.addMiddleware(sch, v)
		if err != nil {
			return err
		}
		err = g.addSecurity(sch, v)
		if err != nil {
			return err
		}
		err = g.addOperation(basePath, sch, v, chain)
		if err != nil {
			return err
		}
	}
	if t.Kind() == gotype.Ptr {
		t = t.Elem()
	}
	if t.Kind() == gotype.Struct {
		numf := t.NumField()
		for i := 0; i != numf; i++ {
			v := t.Field(i)
			if v.IsAnonymous() {
				doc := strings.TrimSpace(v.Doc().Text())
				if doc == "" {
					continue
				}
				_, tag := utils.GetTag(doc)
				basePath := path.Join(basePath, tag.Get("path"))

				v = v.Elem()

				err = g.addMethods(basePath, sch, v, chain, filter)
				if err != nil {
					return err
				}

			} else {
				name := v.Name()
				if !IsExported(name) {
					continue
				}
				doc := strings.TrimSpace(v.Doc().Text())
				if doc == "" {
					continue
				}

				_, tag := utils.GetTag(doc)
				lpath := tag.Get("path")
				if lpath == "" {
					continue
				}
				basePath := path.Join(basePath, lpath)

				v = v.Elem()

				newChain := make([]string, len(chain), len(chain)+1)
				copy(newChain, chain)
				newChain = append(newChain, name)
				filter := map[string]bool{}
				err = g.addMethods(basePath, sch, v, newChain, filter)
				if err != nil {
					return err
				}
			}
		}
	}
	return
}

func (g *Parser) addMiddleware(sch *spec.Type, t gotype.Type) (err error) {
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

	doc, tag := utils.GetTag(doc)
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

	reqs, err := g.addRequests(path, t)
	midd.Requests = reqs

	resps, err := g.addResponses(t)
	midd.Responses = resps

	g.api.Middlewares[key] = midd
	return nil
}

func (g *Parser) addWrapping(sch *spec.Type, t gotype.Type) (err error) {
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

	doc, tag := utils.GetTag(doc)
	name := GetName(oname, tag)
	wrapping := tag.Get("wrapping")
	if wrapping == "" {
		return nil
	}

	path := ""
	route := tag.Get("route")
	if route != "" {
		_, _, path, _ = GetRoute(route)
	}

	hash := utils.Hash(oname, pkgpath, wrapping, doc)
	key := g.namedMidd.GetName(name, hash)

	wrap := &spec.Wrapping{}
	wrap.Name = name
	wrap.PkgPath = pkgpath
	wrap.Schema = wrapping
	wrap.Description = doc
	wrap.Type = sch

	reqs, err := g.addRequests(path, t)
	wrap.Requests = reqs

	resps, err := g.addResponses(t)
	wrap.Responses = resps

	g.api.Wrappings[key] = wrap
	return nil
}

func (g *Parser) addSecurity(sch *spec.Type, t gotype.Type) (err error) {
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

	doc, tag := utils.GetTag(doc)
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

	reqs, err := g.addRequests("", t)
	secu.Requests = reqs

	resps, err := g.addResponses(t)
	secu.Responses = resps

	g.api.Securitys[key] = secu
	return nil
}

func (g *Parser) addOperation(basePath string, sch *spec.Type, t gotype.Type, chain []string) (err error) {
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
	doc, tag := utils.GetTag(doc)
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
	oper.Chain = chain

	if sch != nil {
		sch := sch
		if sch.Ref != "" {
			sch = g.api.Types[sch.Ref]
		}
		oper.Tags = append(oper.Tags, sch.Name)
	}

	reqs, err := g.addRequests(pat, t)
	if err != nil {
		return err
	}
	oper.Requests = reqs

	resps, err := g.addResponses(t)
	if err != nil {
		return err
	}
	oper.Responses = resps

	g.api.Operations = append(g.api.Operations, oper)
	return nil
}

func (g *Parser) addResponses(t gotype.Type) (resps []*spec.Response, err error) {
	numout := t.NumOut()
	for i := 0; i != numout; i++ {
		v := t.Out(i)
		resp, err := g.addResponse(v)
		if err != nil {
			return nil, err
		}
		resps = append(resps, resp)
	}
	return resps, nil
}

func (g *Parser) addResponse(t gotype.Type) (resp *spec.Response, err error) {

	oname := t.Name()
	doc := strings.TrimSpace(t.Comment().Text())
	doc, tag := utils.GetTag(doc)
	name := GetName(oname, tag)
	code := tag.Get("code")
	in := tag.Get("in")
	content := tag.Get("content")
	t = t.Declaration()

	kind := t.Kind()
	if in == "" {
		in = "body"
	}

	sch, err := g.addType(t)
	if err != nil {
		return nil, err
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
			typ := sch
			if typ.Ref != "" {
				typ = g.api.Types[typ.Ref]
			}

			if typ.Attr.Has(spec.AttrImage) {
				content = "image"
			} else if typ.Attr.Has(spec.AttrReader) {
				content = "file"
			} else if kind != gotype.Error {
				content = "json"
			} else {
				content = "error"
			}
		}
	}

	si := sch.Ref
	if si == "" {
		si = sch.Ident
	}
	hash := utils.Hash(oname, in, code, content, doc, si)
	key := g.namedResp.GetName(name+"_"+in, hash)

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

func (g *Parser) addRequests(basePath string, t gotype.Type) (reqs []*spec.Request, err error) {
	numin := t.NumIn()

	for i := 0; i != numin; i++ {
		v := t.In(i)
		req, err := g.addRequest(basePath, v)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (g *Parser) addRequest(basePath string, t gotype.Type) (par *spec.Request, err error) {

	oname := t.Name()
	doc := strings.TrimSpace(t.Comment().Text())
	doc, tag := utils.GetTag(doc)
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
		} else if req, ok := g.importChild("net/url", "Userinfo"); ok && gotype.Equal(t.Elem(), req) {
			return &spec.Request{
				In:    "none",
				Ident: "*net/url.Userinfo",
				Name:  "r.URL.User",
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

	sch, err := g.addType(t)
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
			if basePath != "" && strings.Index(basePath, "{"+name+"}") != -1 {
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
			case spec.Map, spec.Struct:
				in = "body"
			case spec.Interface:
				in = "middleware"
			default:
				if basePath != "" && strings.Index(basePath, "{"+name+"}") != -1 {
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
	key := g.namedReq.GetName(name+"_"+in, hash)

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

func (g *Parser) addType(t gotype.Type) (sch *spec.Type, err error) {
	oname := t.Name()
	pkgpath := t.PkgPath()
	doc := strings.TrimSpace(t.Doc().Text())
	doc, tag := utils.GetTag(doc)
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
		sch.Attr.Add(spec.AttrTextUnmarshaler)
		sch.Attr.Add(spec.AttrTextMarshaler)
		sch.Attr.Add(spec.AttrJSONUnmarshaler)
		sch.Attr.Add(spec.AttrJSONMarshaler)
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
				val, err := g.addType(v.Elem())
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
				vname := v.Name()
				v = v.Declaration()
				if typname := v.Name(); name == typname {
					if value := v.Value(); value != "" {
						desc := v.Doc().Text() + "\n" + v.Comment().Text()
						desc = strings.TrimSpace(desc)
						sch.Enum = append(sch.Enum, &spec.Enum{
							Name:        vname,
							Value:       value,
							Description: desc,
						})
					}
				}
			}
		}
	case gotype.Map:
		schk, err := g.addType(t.Key())
		if err != nil {
			return nil, err
		}
		schv, err := g.addType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch.Key = schk
		sch.Elem = schv
	case gotype.Slice:
		schv, err := g.addType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch.Elem = schv
	case gotype.Array:
		schv, err := g.addType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch.Elem = schv
		sch.Len = t.Len()
	case gotype.Ptr:
		schv, err := g.addType(t.Elem())
		if err != nil {
			return nil, err
		}
		sch.Elem = schv
	case gotype.Interface:
		// No action
	default:
		return nil, fmt.Errorf("Gen.addType: unsupported type: %s is %s kind\n", t.String(), t.Kind().String())
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
