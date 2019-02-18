package openapi

import (
	"fmt"
	"math"
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
		},
	}
}

func (g *GenOpenAPI) SetInfo(info *oaspec.Info) *GenOpenAPI {
	g.openapi.Info = info
	return g
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

	if g.openapi.Info == nil {
		g.openapi.Info = &oaspec.Info{
			Title:       "OpenAPI Demo",
			Description: "Automatically generated",
			Version:     "0.0.1",
			Contact: &oaspec.Contact{
				Name: "wzshiming",
				URL:  "https://github.com/wzshiming/gen",
			},
		}
	}

	return g.openapi, nil
}

func (g *GenOpenAPI) Components() (err error) {

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
		case "header", "cookie", "path", "query":
			par, err := g.Parameters(v)
			if err != nil {
				return err
			}
			g.openapi.Components.Parameters[k] = par
		default:
			// No action
		}
	}

	for k, v := range g.api.Responses {
		switch v.In {
		case "cookie":
			// No action
		case "header":
			// No action
		case "body":
			_, resp, err := g.ResponsesBody(v)
			if err != nil {
				return err
			}
			g.openapi.Components.Responses[k] = resp
		default:
			// No action
		}
	}

	tmpTags := map[*spec.Type]bool{}
	for _, v := range g.api.Operations {
		err := g.Operations(v)
		if err != nil {
			return err
		}
		if v.Type != nil && !tmpTags[v.Type] {
			err := g.Tags(v)
			if err != nil {
				return err
			}
			tmpTags[v.Type] = true
		}
	}

	return nil
}

func (g *GenOpenAPI) Tags(ope *spec.Operation) (err error) {
	typ := ope.Type
	typ = g.api.Types[typ.Ref]
	sch, err := g.Schemas(typ)
	if err != nil {
		return err
	}
	if sch.Ref != "" {
		sch = g.openapi.Components.Schemas[sch.Ref]
	}
	description := sch.Description
	for _, v := range ope.Tags {
		g.openapi.Tags = append(g.openapi.Tags, &oaspec.Tag{
			Name:        v,
			Description: description,
		})
	}
	return nil
}

func (g *GenOpenAPI) Requests(oper *oaspec.Operation, reqs []*spec.Request) (err error) {

	for _, v := range reqs {
		req := v
		if v.Ref != "" {
			req = g.api.Requests[v.Ref]
		}
		switch req.In {
		case "middleware":
			for _, v := range g.api.Middlewares {
				if len(v.Responses) == 0 {
					continue
				}
				resp := v.Responses[0]
				if resp.Ref != "" {
					resp = g.api.Responses[resp.Ref]
				}

				if req.Name == resp.Name {
					err := g.Requests(oper, v.Requests)
					if err != nil {
						return err
					}
				}
			}
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
		case "header", "cookie", "path", "query":
			par, err := g.Parameters(v)
			if err != nil {
				return err
			}
			oper.Parameters = append(oper.Parameters, par)
		default:
			// No action
		}
	}

	return nil
}

func (g *GenOpenAPI) Responses(oper *oaspec.Operation, resps []*spec.Response) (err error) {
	oper.Responses = map[string]*oaspec.Response{}
	headers := map[string]*oaspec.Header{}
	for _, resp := range resps {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		switch resp.In {
		case "cookie":
			// TODO: Process the returned cookie
		case "header":
			name, head, err := g.ResponsesHeader(resp)
			if err != nil {
				return err
			}
			headers[name] = head
		case "body":
			code, resp, err := g.ResponsesBody(resp)
			if err != nil {
				return err
			}

			if len(headers) != 0 {
				resp.Headers = headers
			}

			oper.Responses[code] = resp
		default:
			// No action
		}
	}

	return nil
}

func (g *GenOpenAPI) Operations(ope *spec.Operation) (err error) {
	oper := &oaspec.Operation{}

	err = g.Requests(oper, ope.Requests)
	if err != nil {
		return err
	}

	err = g.Responses(oper, ope.Responses)
	if err != nil {
		return err
	}

	oper.Description = ope.Description
	oper.Summary = ope.Summary
	oper.Deprecated = ope.Deprecated

	if g.openapi.Paths[ope.Path] == nil {
		g.openapi.Paths[ope.Path] = &oaspec.PathItem{}
	}
	item := g.openapi.Paths[ope.Path]

	for _, method := range strings.Split(ope.Method, ",") {
		switch strings.ToLower(method) {
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
	}

	for _, v := range ope.Tags {
		oper.Tags = append(oper.Tags, v)
	}

	return nil
}

func (g *GenOpenAPI) ResponsesHeader(res *spec.Response) (name string, head *oaspec.Header, err error) {
	sch, err := g.Schemas(res.Type)
	if err != nil {
		return "", nil, err
	}
	head = &oaspec.Header{}
	head.Schema = sch
	return res.Name, head, nil
}

func (g *GenOpenAPI) ResponsesBody(res *spec.Response) (code string, resp *oaspec.Response, err error) {
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
		resp = oaspec.NewResponse(res.Content, nil)
	}
	resp.Description = res.Description
	if resp.Description == "" {
		resp.Description = "Response code is " + res.Code
	}
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

	typ := req.Type
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	for _, v := range typ.Enum {
		sch.Enum = append(sch.Enum, oaspec.Any(v.Value))
		if v.Description != "" {
			par.Description += "\n" + v.Value + ":" + v.Description
		}
	}
	par.Description = strings.TrimSpace(par.Description)

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
	case "urlencoded":
		body = oaspec.URLEncodedRequestBody(sch)
	case "formdata":
		body = oaspec.FormDataRequestBody(sch)
	case "octetstream", "file":
		body = oaspec.OctetStreamRequestBody(sch)

		prop := &oaspec.Schema{}
		prop.Type = "object"
		prop.Properties = map[string]*oaspec.Schema{
			req.Name: sch,
		}
		body.Content[oaspec.MimeFormData] = &oaspec.MediaType{
			Schema: prop,
		}

	case "image":
		body = &oaspec.RequestBody{}
		body.Content = map[string]*oaspec.MediaType{
			"image/*": &oaspec.MediaType{
				Schema: sch,
			},
		}

		prop := &oaspec.Schema{}
		prop.Type = "object"
		prop.Properties = map[string]*oaspec.Schema{
			req.Name: sch,
		}
		body.Content[oaspec.MimeFormData] = &oaspec.MediaType{
			Schema: prop,
		}

	default:
		body = oaspec.NewRequestBody(req.Content, sch)
	}
	body.Description = req.Description
	return
}

func (g *GenOpenAPI) Schemas(typ *spec.Type) (sch *oaspec.Schema, err error) {
	if typ.Ref != "" {
		if g.openapi.Components.Schemas[typ.Ref] == nil {
			typ0, err := g.Schemas(g.api.Types[typ.Ref])
			if err != nil {
				return nil, err
			}
			g.openapi.Components.Schemas[typ.Ref] = typ0
		}
		return oaspec.RefSchemas(typ.Ref), nil
	}

	if typ.Attr.HasOne(spec.AttrReader | spec.AttrImage) {
		sch = &oaspec.Schema{}
		sch.Type = "string"
		sch.Format = "binary"
	} else if typ.Attr.Has(spec.AttrTextMarshaler) {
		sch = oaspec.StrFmtProperty(typ.Name)
	} else {
		switch typ.Kind {
		default:
			sch = oaspec.StrFmtProperty(strings.ToLower(typ.Kind.String()))
		case spec.Time:
			sch = oaspec.DateTimeProperty()
		case spec.String:
			sch = oaspec.StringProperty()
		case spec.Bool:
			sch = oaspec.BooleanProperty()
		case spec.Float32:
			sch = oaspec.Float32Property() // .WithMinimum(-math.MaxFloat32, false).WithMaximum(math.MaxFloat32, false)
		case spec.Float64:
			sch = oaspec.Float64Property() // .WithMinimum(-math.MaxFloat64, false).WithMaximum(math.MaxFloat64, false)
		case spec.Int8:
			sch = oaspec.Int8Property().WithMinimum(math.MinInt8, false).WithMaximum(math.MaxInt8, false)
		case spec.Int16:
			sch = oaspec.Int16Property().WithMinimum(math.MinInt16, false).WithMaximum(math.MaxInt16, false)
		case spec.Int32:
			sch = oaspec.Int32Property() // .WithMinimum(math.MinInt32, false).WithMaximum(math.MaxInt32, false)
		case spec.Int64, spec.Int:
			sch = oaspec.Int64Property() // .WithMinimum(math.MinInt64, false).WithMaximum(math.MaxInt64, false)
		case spec.Uint8:
			sch = oaspec.IntFmtProperty("uin8").WithMinimum(0, false).WithMaximum(math.MaxUint8, false)
		case spec.Uint16:
			sch = oaspec.IntFmtProperty("uin16").WithMinimum(0, false).WithMaximum(math.MaxUint16, false)
		case spec.Uint32:
			sch = oaspec.IntFmtProperty("uin32") // .WithMinimum(0, false).WithMaximum(math.MaxUint32, false)
		case spec.Uint64, spec.Uint:
			sch = oaspec.IntFmtProperty("uin64") // .WithMinimum(0, false).WithMaximum(math.MaxUint64, false)
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
				if v.Description != "" {
					val.Description += "\n" + v.Description
				}
				val.Description = strings.TrimSpace(val.Description)

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
	}

	sch.Description = typ.Description
	for _, v := range typ.Enum {
		sch.Enum = append(sch.Enum, oaspec.Any(v.Value))
		if v.Description != "" {
			sch.Description += "\n" + v.Value + ":" + v.Description
		}
	}
	sch.Description = strings.TrimSpace(sch.Description)
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
	secu.Description = sec.Description
	g.openapi.Components.SecuritySchemes[sec.Name] = secu
	return
}
