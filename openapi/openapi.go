package openapi

import (
	"fmt"
	"strings"

	"github.com/wzshiming/gen/spec"
	oaspec "github.com/wzshiming/openapi/spec"
)

type GenOpenAPI struct {
	api     *spec.API
	openapi *oaspec.OpenAPI
	servers []string
}

func NewGenOpenAPI(api *spec.API) *GenOpenAPI {
	return &GenOpenAPI{
		api: api,
		openapi: &oaspec.OpenAPI{
			OpenAPI:    "3.0.1",
			Components: oaspec.NewComponents(),
			Paths:      oaspec.Paths{},
			Info: &oaspec.Info{
				Title:       "OpenAPI Demo",
				Description: "Demo of github.com/wzshiming/openapi",
				Contact: &oaspec.Contact{
					Name:  "wzshiming",
					URL:   "https://github.com/wzshiming",
					Email: "wzshiming@foxmail.com",
				},
				License: &oaspec.License{
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

func (g *GenOpenAPI) Generate() (*oaspec.OpenAPI, error) {
	err := g.Components()
	if err != nil {
		return nil, err
	}
	servers, err := oaspec.NewServers(g.servers...)
	if err != nil {
		return nil, err
	}
	g.openapi.Servers = servers
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

	for _, v := range g.api.Securitys {
		err := g.SecurityScheme(v)
		if err != nil {
			return err
		}
	}

	for k, v := range g.api.Requests {
		switch v.In {
		case "security":
			// No action
		case "body":
			body, err := g.RequestBody(v)
			if err != nil {
				return err
			}
			g.openapi.Components.RequestBodies[k] = body
		default:
			par, err := g.Parameters(v)
			if err != nil {
				return err
			}
			g.openapi.Components.Parameters[k] = par
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

func (g *GenOpenAPI) Operations(ope *spec.Operation) (err error) {
	oper := &oaspec.Operation{}

	for _, v := range ope.Requests {
		req := v
		if v.Ref != "" {
			req = g.api.Requests[v.Ref]
		}
		switch req.In {
		case "security":
			for _, v := range g.api.Securitys {
				if len(v.Responses) == 0 {
					continue
				}
				resp := v.Responses[0]
				if resp.Ref != "" {
					resp = g.api.Responses[resp.Ref]
				}

				if req.Name == resp.Name {
					oper.Security = append(oper.Security, map[string]oaspec.SecurityRequirement{
						v.Name: oaspec.SecurityRequirement{},
					})
				}
			}
		case "body":
			body, err := g.RequestBody(v)
			if err != nil {
				return err
			}
			oper.RequestBody = body
		default:
			par, err := g.Parameters(v)
			if err != nil {
				return err
			}
			oper.Parameters = append(oper.Parameters, par)
		}
	}

	oper.Responses = map[string]*oaspec.Response{}
	for _, v := range ope.Responses {
		code, resp, err := g.Responses(v)
		if err != nil {
			return err
		}
		oper.Responses[code] = resp
	}
	oper.Description = ope.Description

	if g.openapi.Paths[ope.Path] == nil {
		g.openapi.Paths[ope.Path] = &oaspec.PathItem{}
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
			g.openapi.Tags = append(g.openapi.Tags, &oaspec.Tag{
				Name:        v,
				Description: description,
			})
			oper.Tags = append(oper.Tags, v)
		}
	}

	return
}

func (g *GenOpenAPI) Responses(res *spec.Response) (code string, resp *oaspec.Response, err error) {
	if res.Ref != "" {
		return g.api.Responses[res.Ref].Code, oaspec.RefResponse(res.Ref), nil
	}
	sch, err := g.Schemas(res.Type)
	if err != nil {
		return "", nil, err
	}
	switch res.Content {
	case "json":
		resp = oaspec.JSONResponse(sch)
	case "xml":
		resp = oaspec.XMLResponse(sch)
	case "octetstream":
		resp = oaspec.OctetStreamResponse(sch)
	case "urlencoded":
		resp = oaspec.URLEncodedResponse(sch)
	case "formdata":
		resp = oaspec.FormDataResponse(sch)
	case "textplain", "error":
		resp = oaspec.TextPlainResponse(sch)
	default:
		return "", nil, fmt.Errorf("Responses undefined content:%s", res.Content)
	}
	resp.Description = res.Description
	code = res.Code
	return
}

func (g *GenOpenAPI) Parameters(req *spec.Request) (par *oaspec.Parameter, err error) {
	if req.Ref != "" {
		return oaspec.RefParameter(req.Ref), nil
	}
	sch, err := g.Schemas(req.Type)
	if err != nil {
		return nil, err
	}
	switch req.In {
	case "header":
		par = oaspec.HeaderParam(req.Name, sch)
	case "cookie":
		par = oaspec.CookieParam(req.Name, sch)
	case "path":
		par = oaspec.PathParam(req.Name, sch)
	case "query":
		par = oaspec.QueryParam(req.Name, sch)
	default:
		return nil, fmt.Errorf("Parameters undefined in:%s", req.In)
	}
	par.Description = req.Description
	return
}

func (g *GenOpenAPI) RequestBody(req *spec.Request) (body *oaspec.RequestBody, err error) {
	if req.Ref != "" {
		return oaspec.RefRequestBody(req.Ref), nil
	}
	sch, err := g.Schemas(req.Type)
	if err != nil {
		return nil, err
	}

	switch req.Content {
	case "json":
		body = oaspec.JSONRequestBody(sch)
	case "xml":
		body = oaspec.XMLRequestBody(sch)
	case "textplain":
		body = oaspec.TextPlainRequestBody(sch)
	case "octetstream":
		body = oaspec.OctetStreamRequestBody(sch)
	case "urlencoded":
		body = oaspec.URLEncodedRequestBody(sch)
	case "formdata":
		body = oaspec.FormDataRequestBody(sch)
	default:
		return nil, fmt.Errorf("RequestBody undefined content:%s", req.Content)
	}
	body.Description = req.Description
	return
}

func (g *GenOpenAPI) Schemas(typ *spec.Type) (sch *oaspec.Schema, err error) {
	if typ.Ref != "" {
		return oaspec.RefSchemas(typ.Ref), nil
	}
	switch typ.Kind {
	default:
		sch = oaspec.StrFmtProperty(strings.ToLower(typ.Kind.String()))
	case spec.String:
		sch = oaspec.StringProperty()
	case spec.Bool:
		sch = oaspec.BooleanProperty()
	case spec.Float32:
		sch = oaspec.Float32Property() //.WithMinimum(math.SmallestNonzeroFloat32, false).WithMaximum(math.MaxFloat32, false)
	case spec.Float64:
		sch = oaspec.Float64Property() //.WithMinimum(math.SmallestNonzeroFloat64, false).WithMaximum(math.MaxFloat64, false)
	case spec.Int8:
		sch = oaspec.Int8Property() //.WithMinimum(math.MinInt8, false).WithMaximum(math.MaxInt8, false)
	case spec.Int16:
		sch = oaspec.Int16Property() //.WithMinimum(math.MinInt16, false).WithMaximum(math.MaxInt16, false)
	case spec.Int32:
		sch = oaspec.Int32Property() //.WithMinimum(math.MinInt32, false).WithMaximum(math.MaxInt32, false)
	case spec.Int64, spec.Int:
		sch = oaspec.Int64Property() //.WithMinimum(math.MinInt64, false).WithMaximum(math.MaxInt64, false)
	case spec.Uint8:
		sch = oaspec.IntFmtProperty("uin8") //.WithMinimum(0, false).WithMaximum(math.MaxUint8, false)
	case spec.Uint16:
		sch = oaspec.IntFmtProperty("uin16") //.WithMinimum(0, false).WithMaximum(math.MaxUint16, false)
	case spec.Uint32:
		sch = oaspec.IntFmtProperty("uin32") //.WithMinimum(0, false).WithMaximum(math.MaxUint32, false)
	case spec.Uint64, spec.Uint:
		sch = oaspec.IntFmtProperty("uin64") //.WithMinimum(0, false).WithMaximum(math.MaxUint64, false)
	case spec.Map:
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
		sch = oaspec.MapProperty(sch)
	case spec.Slice:
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
		sch = oaspec.ArrayProperty(sch)
	case spec.Array:
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
		sch = oaspec.ArrayProperty(sch).WithMaxItems(int64(typ.Len))
	case spec.Ptr:
		sch, err = g.Schemas(typ.Elem)
		if err != nil {
			return nil, err
		}
	case spec.Error:
		sch = oaspec.StrFmtProperty("error")
	case spec.Struct:
		sch = &oaspec.Schema{}
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
				sch.Properties = map[string]*oaspec.Schema{}
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
		sch.Enum = append(sch.Enum, oaspec.Any(v))
	}
	return sch, nil
}

func (g *GenOpenAPI) SecurityScheme(sec *spec.Security) (err error) {
	secu := &oaspec.SecurityScheme{}
	secu.Type = sec.Schema
	if len(sec.Requests) == 0 {
		return nil
	}
	req := sec.Requests[0]
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}
	secu.In = req.In
	secu.Name = req.Name
	secu.Description = secu.Description
	g.openapi.Components.SecuritySchemes[sec.Name] = secu
	return
}
