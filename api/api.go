package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/wzshiming/gotype"
	"github.com/wzshiming/openapi/spec"
)

// GenAPI is the parse type generating definitions
type GenAPI struct {
	imp *gotype.Importer
	spec.Components
	paths map[string]*spec.PathItem
}

func NewGenAPI() *GenAPI {
	return &GenAPI{
		imp:        gotype.NewImporter(gotype.WithCommentLocator()),
		Components: *spec.NewComponents(),
		paths:      map[string]*spec.PathItem{},
	}
}

func (g *GenAPI) Import(pkgpath string) error {
	pkg, err := g.imp.Import(pkgpath)
	if err != nil {
		return err
	}

	numchi := pkg.NumChild()

	for i := 0; i != numchi; i++ {
		v := pkg.Child(i)
		switch v.Kind() {
		case gotype.Func:
			_, err := g.AddPathItem("", v)
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

func (g *GenAPI) OpenAPI() (swa *spec.OpenAPI) {
	swa = &spec.OpenAPI{}

	swa.OpenAPI = "3.0.1"

	swa.Info = &spec.Info{
		Title:       "OpenAPI Demo",
		Description: "Demo of github.com/wzshiming/openapi",
		//TermsOfService: "https://github.com/wzshiming/openapi",
		Contact: &spec.Contact{
			Name:  "wzshiming",
			URL:   "https://github.com/wzshiming",
			Email: "wzshiming@foxmail.com",
		},
		License: &spec.License{
			Name: "MIT",
			URL:  "https://github.com/wzshiming/openapi/blob/master/LICENSE",
		},
		Version: "0.0.1",
	}
	swa.Tags = append(swa.Tags, &spec.Tag{
		Name:        "default",
		Description: "Description",
	})

	swa.Servers, _ = spec.NewServers("http://127.0.0.1:8080/v1/")
	swa.Paths = g.paths
	swa.Components = &g.Components
	return
}

var key = regexp.MustCompile(`#[^#]+#`)

func GetTag(text string) reflect.StructTag {
	dd := key.FindAllString(text, -1)
	for i := 0; i != len(dd); i++ {
		dd[i] = dd[i][1 : len(dd[i])-1]
	}
	return reflect.StructTag(strings.Join(dd, " "))
}

func (g *GenAPI) AddPaths(t gotype.Type) (err error) {
	numm := t.NumMethods()
	if numm == 0 {
		return nil
	}
	tag := GetTag(t.Doc().Text())
	route := tag.Get("route")
	for i := 0; i != numm; i++ {
		v := t.Methods(i)
		_, err := g.AddPathItem(route, v)
		if err != nil {
			return err
		}

	}
	return
}

func (g *GenAPI) AddPathItem(bashPath string, t gotype.Type) (item *spec.PathItem, err error) {
	if t.Kind() != gotype.Func {
		return nil, fmt.Errorf("GenAPI: unsupported type: %s", t.Kind().String())
	}

	tag := GetTag(t.Doc().Text())
	route := tag.Get("route")
	if route == "" {
		return nil, nil
	}
	rs := strings.SplitN(route, " ", 2)
	if len(rs) != 2 {
		return nil, nil
	}
	pat := strings.TrimSpace(rs[1])
	method := strings.ToUpper(strings.TrimSpace(rs[0]))

	if bashPath != "" {
		pat = path.Join(bashPath, pat)
	}

	if g.paths[pat] == nil {
		g.paths[pat] = &spec.PathItem{}
	}

	oper, err := g.AddOperation(pat, t)
	if err != nil {
		return nil, err
	}

	switch method {
	case "GET":
		g.paths[pat].Get = oper
	case "PUT":
		g.paths[pat].Put = oper
	case "POST":
		g.paths[pat].Post = oper
	case "DELETE":
		g.paths[pat].Delete = oper
	case "OPTIONS":
		g.paths[pat].Options = oper
	case "HEAD":
		g.paths[pat].Head = oper
	case "PATCH":
		g.paths[pat].Patch = oper
	case "TRACE":
		g.paths[pat].Trace = oper
	}
	return g.paths[pat], nil
}

func (g *GenAPI) AddOperation(path string, t gotype.Type) (oper *spec.Operation, err error) {
	if t.Kind() != gotype.Func {
		return nil, fmt.Errorf("GenAPI: unsupported type: %s", t.Kind().String())
	}

	oper = &spec.Operation{}
	oper.OperationID = t.Name()
	{
		numin := t.NumIn()
		for i := 0; i != numin; i++ {
			v := t.In(i)
			par, req, err := g.AddParameter(path, v)
			if err != nil {
				return nil, err
			}
			if par != nil {
				oper.Parameters = append(oper.Parameters, par)
			}
			if req != nil {
				oper.RequestBody = req
			}
		}
	}

	{
		numout := t.NumOut()
		for i := 0; i != numout; i++ {
			v := t.Out(i)
			code, resp, err := g.AddResponse(v)
			if err != nil {
				return nil, err
			}
			if oper.Responses == nil {
				oper.Responses = map[string]*spec.Response{}
			}
			for _, code := range strings.Split(code, ",") {
				oper.Responses[code] = resp
			}
		}
	}
	return oper, nil
}

func (g *GenAPI) AddResponse(t gotype.Type) (code string, resp *spec.Response, err error) {
	tag := GetTag(t.Comment().Text())
	name := t.Name()
	code = tag.Get("code")
	content := tag.Get("code")
	if content == "" {
		content = "json"
	}
	if code == "" {
		if t.Elem().Kind() != gotype.Error {
			code = "200"
		} else {
			code = "400"
		}
	}
	sch, err := g.AddSchemas(t.Elem())
	if err != nil {
		return "", nil, err
	}

	key := name + "." + hash(content+sch.Ref)
	if g.Responses[key] != nil {
		return code, spec.RefResponse(key), nil
	}
	defer func() {
		if resp != nil {
			g.Responses[key] = resp
			resp = spec.RefResponse(key)
		}
	}()

	switch content {
	case "json":
		resp = spec.JSONResponse(sch)
	case "xml":
		resp = spec.XMLResponse(sch)
	case "urlencoded":
		resp = spec.URLEncodedResponse(sch)
	case "formdata":
		resp = spec.FormDataResponse(sch)
	}

	return code, resp, nil
}

func (g *GenAPI) AddParameter(path string, t gotype.Type) (par *spec.Parameter, req *spec.RequestBody, err error) {
	rawname := t.Name()
	tag := GetTag(t.Comment().Text())

	in := tag.Get("in")
	if in == "" {
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

	content := tag.Get("code")
	if content == "" {
		content = "json"
	}

	name, ok := tag.Lookup("name")
	if !ok {
		name = rawname
	}

	sch, err := g.AddSchemas(t.Elem())
	if err != nil {
		return nil, nil, err
	}

	key := name + "." + hash(name+in+sch.Ref)
	bodyKey := name + "." + content + "." + hash(name+content+sch.Ref)
	if in == "body" {
		if g.RequestBodies[bodyKey] != nil {
			return nil, spec.RefRequestBody(bodyKey), nil
		}
	} else {
		if g.Parameters[key] != nil {
			return spec.RefParameter(key), nil, nil
		}
	}

	defer func() {
		if par != nil {
			g.Parameters[key] = par
			par = spec.RefParameter(key)
		} else if req != nil {
			g.RequestBodies[bodyKey] = req
			req = spec.RefRequestBody(bodyKey)
		}

	}()

	switch in {
	case "header":
		par = spec.HeaderParam(name, sch)
	case "path":
		par = spec.PathParam(name, sch)
	case "query":
		par = spec.QueryParam(name, sch)
	case "cookie":
		par = spec.CookieParam(name, sch)
	case "body":
		switch content {
		case "json":
			req = spec.JSONRequestBody(sch)
		case "xml":
			req = spec.XMLRequestBody(sch)
		case "urlencoded":
			req = spec.URLEncodedRequestBody(sch)
		case "formdata":
			req = spec.FormDataRequestBody(sch)
		default:
			return nil, nil, fmt.Errorf("undefined content:%s", content)
		}
	default:
		return nil, nil, fmt.Errorf("undefined in:%s", in)
	}
	return par, req, nil
}

func (g *GenAPI) AddSchemas(t gotype.Type) (sch *spec.Schema, err error) {
	typname := t.Name()
	pkgpath := t.PkgPath()
	key := typname + "." + hash(pkgpath)
	if g.Schemas[key] != nil {
		return spec.RefSchemas(key), nil
	}
	if t.IsGoroot() && pkgpath == "time" && typname == "Time" {
		return spec.DateTimeProperty(), nil
	}

	kind := t.Kind()
	if typname != "" && typname != strings.ToLower(kind.String()) {
		defer func() {
			if sch == nil {
				return
			}
			g.Schemas[key] = sch
			sch = spec.RefSchemas(key)
		}()
	}

	switch kind {
	case gotype.Struct:
		anons := []*spec.Schema{}
		// Anonymo field
		{
			num := t.NumAnonymo()
			for i := 0; i != num; i++ {
				v := t.Anonymo(i)
				val, err := g.AddSchemas(v.Elem())
				if err != nil {
					return nil, err
				}
				anons = append(anons, val)
			}
		}

		sch = &spec.Schema{}
		// Field
		{
			num := t.NumField()
			for i := 0; i != num; i++ {
				v := t.Field(i)
				name := v.Name()
				tag := v.Tag()

				val, err := g.AddSchemas(v.Elem())
				if err != nil {
					return nil, err
				}
				if sch.Properties == nil {
					sch.Properties = map[string]*spec.Schema{}
				}

				if n := strings.Split(tag.Get("json"), ",")[0]; n != "" {
					if n == "-" {
						continue
					}
					name = n
				}
				if n := strings.Split(tag.Get("xml"), ",")[0]; n != "" {
					if n == "-" {
						continue
					}
					val = val.WithXMLName(n)
				}
				sch.Properties[name] = val
			}
		}

		if len(anons) == 0 {
			return sch, nil
		}
		anons = append(anons, sch)
		sch = spec.ComposedSchema(anons...)
	case gotype.String:
		sch = spec.StringProperty()
	case gotype.Bool:
		sch = spec.BooleanProperty()
	case gotype.Float32:
		sch = spec.Float32Property().WithMinimum(math.SmallestNonzeroFloat32, false).WithMaximum(math.MaxFloat32, false)
	case gotype.Float64:
		sch = spec.Float64Property().WithMinimum(math.SmallestNonzeroFloat64, false).WithMaximum(math.MaxFloat64, false)
	case gotype.Int8:
		sch = spec.Int8Property().WithMinimum(math.MinInt8, false).WithMaximum(math.MaxInt8, false)
	case gotype.Int16:
		sch = spec.Int16Property().WithMinimum(math.MinInt16, false).WithMaximum(math.MaxInt16, false)
	case gotype.Int32:
		sch = spec.Int32Property().WithMinimum(math.MinInt32, false).WithMaximum(math.MaxInt32, false)
	case gotype.Int64, gotype.Int:
		sch = spec.Int64Property().WithMinimum(math.MinInt64, false).WithMaximum(math.MaxInt64, false)
	case gotype.Uint8:
		sch = spec.IntFmtProperty("uin8").WithMinimum(0, false).WithMaximum(math.MaxUint8, false)
	case gotype.Uint16:
		sch = spec.IntFmtProperty("uin16").WithMinimum(0, false).WithMaximum(math.MaxUint16, false)
	case gotype.Uint32:
		sch = spec.IntFmtProperty("uin32").WithMinimum(0, false).WithMaximum(math.MaxUint32, false)
	case gotype.Uint64, gotype.Uint:
		sch = spec.IntFmtProperty("uin64").WithMinimum(0, false).WithMaximum(math.MaxUint64, false)
	case gotype.Map:
		sch, err = g.AddSchemas(t.Elem())
		if err != nil {
			return nil, err
		}
		sch = spec.MapProperty(sch)
	case gotype.Slice:
		sch, err = g.AddSchemas(t.Elem())
		if err != nil {
			return nil, err
		}
		sch = spec.ArrayProperty(sch)
	case gotype.Array:
		sch, err = g.AddSchemas(t.Elem())
		if err != nil {
			return nil, err
		}
		sch = spec.ArrayProperty(sch).WithMaxItems(int64(t.Len()))
	case gotype.Ptr:
		sch, err = g.AddSchemas(t.Elem())
		if err != nil {
			return nil, err
		}
	case gotype.Error:
		sch = spec.StrFmtProperty("error")
	default:
		return nil, fmt.Errorf("gotype: unsupported type: %s", t.Kind().String())
	}

	switch kind {
	case gotype.String,
		gotype.Int, gotype.Int8, gotype.Int16, gotype.Int32, gotype.Int64,
		gotype.Uint, gotype.Uint8, gotype.Uint16, gotype.Uint32, gotype.Uint64:
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
				value := v.Value()
				if name == "_" {
					continue
				}
				sch.Enum = append(sch.Enum, spec.Any(value))
			}
		}
	}

	tag := GetTag(t.Doc().Text())
	if format := tag.Get("format"); format != "" {
		if typ := tag.Get("type"); typ != "" {
			sch.WithType(typ, format)
		} else {
			sch.WithFormat(format)
		}
	}

	return sch, nil
}

func hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
