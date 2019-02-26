package client

import (
	"sort"

	"github.com/wzshiming/gen/model"
	"github.com/wzshiming/gen/named"
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/srcgen"
	"github.com/wzshiming/gen/utils"
)

// GenClient is the generating generating
type GenClient struct {
	api *spec.API
	buf *srcgen.File
	model.GenModel
	named *named.Named
}

func NewGenClient(api *spec.API) *GenClient {
	buf := &srcgen.File{}
	return &GenClient{
		api:      api,
		buf:      buf,
		GenModel: *model.NewGenModel(api, buf, ""),
		named:    named.NewNamed("_"),
	}
}

func (g *GenClient) Generate() (*srcgen.File, error) {

	err := g.generateSchemas()
	if err != nil {
		return nil, err
	}
	err = g.generateSecuritys()
	if err != nil {
		return nil, err
	}
	err = g.generateClient()
	if err != nil {
		return nil, err
	}
	return g.buf, nil
}

func (g *GenClient) generateSchemas() (err error) {
	schemas := g.api.Types
	schKey := make([]string, 0, len(schemas))
	for k := range schemas {
		schKey = append(schKey, k)
	}
	sort.Strings(schKey)
	for _, k := range schKey {
		v := schemas[k]
		if v.Attr.Has(spec.AttrRoot) {
			continue
		}
		g.buf.WriteString(utils.CommentLine(v.Description))
		g.buf.WriteString("type ")
		g.buf.WriteString(g.getTypeName(v))
		g.buf.WriteByte(' ')
		err = g.TypesDefine(v)
		if err != nil {
			return err
		}
		g.buf.WriteString("\n\n")
	}
	return
}

func (g *GenClient) generateParameterRequests(req *spec.Request, typ string) (err error) {
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}
	g.buf.WriteFormat("%s ", g.getVarName(req.Name, req.Type))
	if typ != "" {
		g.buf.WriteString(typ)
	} else {
		err = g.Types(req.Type)
	}
	if err != nil {
		return err
	}
	if req.Description != "" {
		g.buf.WriteFormat("/* %s */", utils.MergeLine(req.Description))
	}
	return nil
}

func (g *GenClient) generateParameterResponses(resp *spec.Response) (err error) {
	if resp.Ref != "" {
		resp = g.api.Responses[resp.Ref]
	}
	g.buf.WriteFormat("%s ", g.getVarName(resp.Name, resp.Type))
	err = g.Types(resp.Type)
	if err != nil {
		return err
	}
	if resp.Description != "" {
		g.buf.WriteFormat("/* %s */", utils.MergeLine(resp.Description))
	}
	return nil
}
