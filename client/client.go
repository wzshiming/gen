package client

import (
	"sort"

	"github.com/wzshiming/gen/model"
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/srcgen"
	"github.com/wzshiming/gen/utils"
)

// GenClient is the generating generating
type GenClient struct {
	api *spec.API
	buf *srcgen.File
	model.GenModel
}

func NewGenClient(api *spec.API) *GenClient {
	buf := &srcgen.File{}
	return &GenClient{
		api:      api,
		buf:      buf,
		GenModel: *model.NewGenModel(api, buf),
	}
}

func (g *GenClient) Generate() (*srcgen.File, error) {

	err := g.GenerateSchemas()
	if err != nil {
		return nil, err
	}
	err = g.GenerateSecuritys()
	if err != nil {
		return nil, err
	}
	err = g.GenerateClient()
	if err != nil {
		return nil, err
	}
	return g.buf, nil
}

func (g *GenClient) GenerateSchemas() (err error) {
	schemas := g.api.Types
	ks := make([]string, 0, len(schemas))
	for k, _ := range schemas {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		v := schemas[k]
		g.buf.WriteString(utils.CommentLine(v.Description))
		g.buf.WriteString("type ")
		g.buf.WriteString(utils.GetName(k))
		g.buf.WriteByte(' ')
		err = g.Types(v)
		if err != nil {
			return err
		}
		g.buf.WriteString("\n\n")
	}
	return
}

func (g *GenClient) GenerateParameterRequests(req *spec.Request) (err error) {
	if req.Ref != "" {
		req = g.api.Requests[req.Ref]
	}
	g.buf.WriteFormat("_%s ", req.Name)
	err = g.Types(req.Type)
	if err != nil {
		return err
	}
	if req.Description != "" {
		g.buf.WriteFormat("/* %s */", utils.MergeLine(req.Description))
	}
	return nil
}

func (g *GenClient) GenerateParameterResponses(resp *spec.Response) (err error) {
	if resp.Ref != "" {
		resp = g.api.Responses[resp.Ref]
	}
	g.buf.WriteFormat("_%s ", resp.Name)
	err = g.Types(resp.Type)
	if err != nil {
		return err
	}
	if resp.Description != "" {
		g.buf.WriteFormat("/* %s */", utils.MergeLine(resp.Description))
	}
	return nil
}
