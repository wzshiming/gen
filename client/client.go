package client

import (
	"sort"

	"github.com/wzshiming/gen/model"
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/srcgen"
	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/namecase"
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

func (g *GenClient) Generate() ([]byte, error) {

	g.buf.WithPackname("main")

	err := g.GenerateSchemas()
	if err != nil {
		return nil, err
	}
	err = g.GenerateOperations()
	if err != nil {
		return nil, err
	}
	return g.buf.Bytes(), nil
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

func (g *GenClient) GenerateOperations() (err error) {
	g.buf.AddImport("", "github.com/wzshiming/requests")
	g.buf.WriteFormat(`
	var Client = requests.NewClient().NewRequest()
`)

	operations := g.api.Operations
	for _, v := range operations {
		err = g.Operations(v)
		if err != nil {
			return err
		}
		err = g.GenerateFuncBody(v)
		if err != nil {
			return err
		}
		g.buf.WriteString("\n\n")
	}
	return
}

func (g *GenClient) GenerateFuncBody(oper *spec.Operation) (err error) {

	g.buf.WriteString("{\n")
	defer g.buf.WriteString(`
	return
}
`)
	g.GenerateRequests(oper)
	g.GenerateResponses(oper)
	return nil
}

func (g *GenClient) GenerateResponses(oper *spec.Operation) (err error) {

	g.buf.WriteString(`
	switch code := resp.StatusCode(); code {
`)

	for _, resp := range oper.Responses {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if resp.Code == "" {
			continue
		}
		g.buf.WriteFormat("case %s:", resp.Code)
		switch resp.Content {
		case "json":
			g.buf.AddImport("", "encoding/json")
			g.buf.WriteFormat(`
	err = json.Unmarshal(resp.Body(),&_%s)
`, resp.Name)
		case "xml":
			g.buf.AddImport("", "encoding/xml")
			g.buf.WriteFormat(`
	err = xml.Unmarshal(resp.Body(),&_%s)
`, resp.Name)
		case "error":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`
	_%s = fmt.Errorf(string(resp.Body()))
`, resp.Name)
		}
		// TODO
	}
	g.buf.AddImport("", "net/http")
	g.buf.WriteString(`default:
		err = fmt.Errorf("Undefined code %d %s", code, http.StatusText(code))
	}
`)

	g.GenerateErrror(oper.Responses)
	return nil
}

func (g *GenClient) GenerateErrror(resps []*spec.Response) (err error) {
	g.buf.WriteString(`
if err != nil {
	return `)
	for i, resp := range resps {
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if i != len(resps)-1 {
			g.TypesZero(resp.Type)
		} else {
			g.buf.WriteString("err")
		}
	}
	g.buf.WriteString(`	
}
`)
	return
}

func (g *GenClient) GenerateRequests(oper *spec.Operation) (err error) {

	g.buf.WriteString("resp, err := Client.Clone().\n")
	for _, v := range oper.Requests {
		req := v
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		switch req.In {
		case "header":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`SetHead("%s", fmt.Sprint(_%s)).
`, req.Name, req.Name)
		case "cookie":
			// TODO
		case "path":
			g.buf.AddImport("", "fmt")
			g.buf.WriteFormat(`SetPath("%s", fmt.Sprint(_%s)).
`, req.Name, req.Name)
		case "query":
			g.buf.WriteFormat(`SetQuery("%s", fmt.Sprint(_%s)).
`, req.Name, req.Name)
		case "body":
			switch req.Content {
			case "json":
				g.buf.WriteFormat("SetJSON(_%s).\n", req.Name)
			case "xml":
				g.buf.WriteFormat("SetXML(_%s).\n", req.Name)
			}
		}
	}
	g.buf.WriteString(namecase.ToPascal(oper.Method) + "(\"" + oper.Path + "\")\n")

	g.GenerateErrror(oper.Responses)
	return nil
}

func (g *GenClient) Operations(oper *spec.Operation) (err error) {

	g.buf.WriteString(utils.CommentLine(oper.Description))
	g.buf.WriteString("func ")

	if oper.Type != nil {
		g.buf.WriteByte('(')
		err = g.Types(oper.Type)
		if err != nil {
			return err
		}
		g.buf.WriteByte(')')
	}
	g.buf.WriteString(utils.GetName(oper.Name))
	g.buf.WriteByte('(')
	for i, v := range oper.Requests {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		err = g.GenerateParameterRequests(v)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(")(")

	for i, v := range oper.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		err = g.GenerateParameterResponses(v)
		if err != nil {
			return err
		}
	}

	g.buf.WriteByte(')')
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
