package openapi

import (
	"fmt"
	"strings"

	"github.com/wzshiming/gen"
	"github.com/wzshiming/openapi/spec"
)

type GenOpenAPI struct {
	api     *gen.API
	openapi *spec.OpenAPI
	servers []string
}

func NewGenOpenAPI(api *gen.API) *GenOpenAPI {
	return &GenOpenAPI{
		api: api,
		openapi: &spec.OpenAPI{
			OpenAPI:    "3.0.1",
			Components: spec.NewComponents(),
			Paths:      spec.Paths{},
			Info: &spec.Info{
				Title:       "OpenAPI Demo",
				Description: "Demo of github.com/wzshiming/openapi",
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
			},
		},
	}
}

func (g *GenOpenAPI) WithServices(servers ...string) *GenOpenAPI {
	g.servers = append(g.servers, servers...)
	return g
}

func (g *GenOpenAPI) Generate() (*spec.OpenAPI, error) {
	err := g.Components()
	if err != nil {
		return nil, err
	}
	servers, err := spec.NewServers(g.servers...)
	if err != nil {
		return nil, err
	}
	g.openapi.Servers = servers
	g.openapi.Security = append(g.openapi.Security, map[string]spec.SecurityRequirement{
		"api_key": spec.SecurityRequirement{},
	})
	sr := &spec.SecurityScheme{}
	sr.In = "query"
	sr.Type = "apiKey"
	sr.Name = "api_key"
	g.openapi.Components.SecuritySchemes["api_key"] = sr
	return g.openapi, nil
}

func (g *GenOpenAPI) Components() (err error) {
	for k, v := range g.api.Types {
		sch, err := g.Schemas(v)
		if err != nil {
			return err
		}
		g.openapi.Components.Schemas[k] = sch

	}

	for k, v := range g.api.Requests {
		par, body, err := g.Parameters(v)
		if err != nil {
			return err
		}
		if par != nil {
			g.openapi.Components.Parameters[k] = par
		} else if body != nil {
			g.openapi.Components.RequestBodies[k] = body
		}
	}

	for k, v := range g.api.Responses {
		_, resp, err := g.Responses(v)
		if err != nil {
			return err
		}
		g.openapi.Components.Responses[k] = resp
	}

	for _, v := range g.api.Operations {
		err := g.Operations(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GenOpenAPI) Operations(ope *gen.Operation) (err error) {
	oper := &spec.Operation{}
	for _, v := range ope.Requests {
		par, body, err := g.Parameters(v)
		if err != nil {
			return err
		}
		if par != nil {
			oper.Parameters = append(oper.Parameters, par)
		} else if body != nil {
			oper.RequestBody = body
		}
	}

	oper.Responses = map[string]*spec.Response{}
	for _, v := range ope.Responses {
		code, resp, err := g.Responses(v)
		if err != nil {
			return err
		}
		oper.Responses[code] = resp
	}
	oper.Description = ope.Description

	if g.openapi.Paths[ope.Path] == nil {
		g.openapi.Paths[ope.Path] = &spec.PathItem{}
	}
	item := g.openapi.Paths[ope.Path]

	switch ope.Method {
	case "get":
		item.Get = oper
	case "put":
		item.Put = oper
	case "post":
		item.Post = oper
	case "delete":
		item.Delete = oper
	case "options":
		item.Options = oper
	case "head":
		item.Head = oper
	case "patch":
		item.Patch = oper
	case "trace":
		item.Trace = oper
	}

	if ope.Type != nil {
		for _, v := range ope.Tags {
			sch, err := g.Schemas(ope.Type)
			if err != nil {
				return err
			}
			description := ""
			if sch.Ref != "" {
				description = g.openapi.Components.Schemas[sch.Ref].Description
			} else {
				description = sch.Description
			}
			g.openapi.Tags = append(g.openapi.Tags, &spec.Tag{
				Name:        v,
				Description: description,
			})
			oper.Tags = append(oper.Tags, v)
		}
	}

	return
}

func (g *GenOpenAPI) Responses(res *gen.Response) (code string, resp *spec.Response, err error) {
	if res.Ref != "" {
		return g.api.Responses[res.Ref].Code, spec.RefResponse(res.Ref), nil
	}
	sch, err := g.Schemas(res.Type)
	if err != nil {
		return "", nil, err
	}
	switch res.Content {
	case "json":
		resp = spec.JSONResponse(sch)
	case "xml":
		resp = spec.XMLResponse(sch)
	case "octetstream":
		resp = spec.OctetStreamResponse(sch)
	case "urlencoded":
		resp = spec.URLEncodedResponse(sch)
	case "formdata":
		resp = spec.FormDataResponse(sch)
	default:
		return "", nil, fmt.Errorf("Responses undefined content:%s", res.Content)
	}
	resp.Description = res.Description
	code = res.Code
	return
}

func (g *GenOpenAPI) Parameters(req *gen.Request) (par *spec.Parameter, body *spec.RequestBody, err error) {
	if req.Ref != "" {
		if _, ok := g.openapi.Components.Parameters[req.Ref]; ok {
			return spec.RefParameter(req.Ref), nil, nil
		} else if _, ok := g.openapi.Components.RequestBodies[req.Ref]; ok {
			return nil, spec.RefRequestBody(req.Ref), nil
		}
		return nil, nil, fmt.Errorf("Responses undefined ref:%s", req.Ref)
	}
	sch, err := g.Schemas(req.Type)
	if err != nil {
		return nil, nil, err
	}
	switch req.In {
	case "header":
		par = spec.HeaderParam(req.Name, sch)
	case "cookie":
		par = spec.CookieParam(req.Name, sch)
	case "path":
		par = spec.PathParam(req.Name, sch)
	case "query":
		par = spec.QueryParam(req.Name, sch)
	case "body":
		switch req.Content {
		case "json":
			body = spec.JSONRequestBody(sch)
		case "xml":
			body = spec.XMLRequestBody(sch)
		case "textplain":
			body = spec.TextPlainRequestBody(sch)
		case "octetstream":
			body = spec.OctetStreamRequestBody(sch)
		case "urlencoded":
			body = spec.URLEncodedRequestBody(sch)
		case "formdata":
			body = spec.FormDataRequestBody(sch)
		default:
			return nil, nil, fmt.Errorf("RequestBody undefined content:%s", req.Content)
		}
	default:
		return nil, nil, fmt.Errorf("Parameters undefined in:%s", req.In)
	}
	if par != nil {
		par.Description = req.Description
	} else if body != nil {
		body.Description = req.Description
	}
	return
}

func (g *GenOpenAPI) Schemas(typ *gen.Type) (sch *spec.Schema, err error) {
	if typ.Ref != "" {
		return spec.RefSchemas(typ.Ref), nil
	}
	switch typ.Type {
	default:
		sch = spec.StrFmtProperty(typ.Type)
	case "string":
		sch = spec.StringProperty()
	case "bool":
		sch = spec.BooleanProperty()
	case "float32":
		sch = spec.Float32Property() //.WithMinimum(math.SmallestNonzeroFloat32, false).WithMaximum(math.MaxFloat32, false)
	case "float64":
		sch = spec.Float64Property() //.WithMinimum(math.SmallestNonzeroFloat64, false).WithMaximum(math.MaxFloat64, false)
	case "int8":
		sch = spec.Int8Property() //.WithMinimum(math.MinInt8, false).WithMaximum(math.MaxInt8, false)
	case "int16":
		sch = spec.Int16Property() //.WithMinimum(math.MinInt16, false).WithMaximum(math.MaxInt16, false)
	case "int32":
		sch = spec.Int32Property() //.WithMinimum(math.MinInt32, false).WithMaximum(math.MaxInt32, false)
	case "int64", "int":
		sch = spec.Int64Property() //.WithMinimum(math.MinInt64, false).WithMaximum(math.MaxInt64, false)
	case "uint8":
		sch = spec.IntFmtProperty("uin8") //.WithMinimum(0, false).WithMaximum(math.MaxUint8, false)
	case "uint16":
		sch = spec.IntFmtProperty("uin16") //.WithMinimum(0, false).WithMaximum(math.MaxUint16, false)
	case "uint32":
		sch = spec.IntFmtProperty("uin32") //.WithMinimum(0, false).WithMaximum(math.MaxUint32, false)
	case "uint64", "uint":
		sch = spec.IntFmtProperty("uin64") //.WithMinimum(0, false).WithMaximum(math.MaxUint64, false)
	case "map":
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
		sch = spec.MapProperty(sch)
	case "slice":
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
		sch = spec.ArrayProperty(sch)
	case "array":
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
		sch = spec.ArrayProperty(sch).WithMaxItems(int64(typ.Len))
	case "ptr":
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
	case "error":
		sch = spec.StrFmtProperty("error")
	case "struct":
		sch = &spec.Schema{}
		sch.Type = "object"
		for _, v := range typ.Fields {
			name := v.Name
			tag := v.Tag

			val, err := g.Schemas(v.Type)
			if err != nil {
				return nil, err
			}
			val.Description += v.Description
			if v.Anonymous {

				sch.AllOf = append(sch.AllOf, val)
				continue
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
	sch.Description = typ.Description
	for _, v := range typ.Enum {
		sch.Enum = append(sch.Enum, spec.Any(v))
	}
	return sch, nil
}
