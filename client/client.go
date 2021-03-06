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
	model.GenModel
	api     *spec.API
	buf     *srcgen.File
	named   *named.Named
	explode bool
}

func NewGenClient(api *spec.API) *GenClient {
	buf := &srcgen.File{}
	return &GenClient{
		api:      api,
		buf:      buf,
		GenModel: *model.NewGenModel(api, buf, api.Imports),
		named:    named.NewNamed("_"),
	}
}

func (g *GenClient) SetExplode(b bool) *GenClient {
	g.explode = b
	return g
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
		if _, ok := g.GenModel.GetPkgPath(v.PkgPath); ok {
			continue
		}
		g.buf.WriteString(utils.CommentLine(v.Description))
		g.buf.WriteString("type ")
		name := g.getTypeName(v)
		g.buf.WriteString(name)
		g.buf.WriteByte(' ')
		err = g.TypesDefine(v)
		if err != nil {
			return err
		}
		if len(v.Enum) != 0 {
			g.buf.WriteString(`
const (
`)
			for _, enum := range v.Enum {
				g.buf.WriteFormat("%s %s = %s\n", g.getEnumName(enum.Name, enum.Value), name, enum.Value)
			}
			g.buf.WriteString(`)`)
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

	desc := req.Description
	_, tag := utils.GetTag(desc)
	if tag.Get("name") == "" {
		desc += "\n#name:\"" + req.Name + "\"#"
	}
	g.buf.WriteFormat("/* %s */", utils.MergeLine(desc))

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

	desc := resp.Description
	_, tag := utils.GetTag(desc)
	if tag.Get("name") == "" {
		desc += "\n#name:\"" + resp.Name + "\"#"
	}
	g.buf.WriteFormat("/* %s */", utils.MergeLine(desc))

	return nil
}
