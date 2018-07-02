package client

import (
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
)

func (g *GenClient) GenerateClient() (err error) {
	g.buf.AddImport("", "github.com/wzshiming/requests")
	g.buf.WriteFormat(`
	var Client = requests.NewClient().NewRequest()
`)

	operations := g.api.Operations
	for _, v := range operations {
		err = g.GenerateOperations(v)
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

func (g *GenClient) GenerateOperations(oper *spec.Operation) (err error) {

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
	reqs := []*spec.Request{}
	for _, req := range oper.Requests {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		switch req.In {
		case "security":
			// No action
		case "header", "path", "query", "body":
			reqs = append(reqs, req)
		}
	}
	for i, req := range reqs {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		err = g.GenerateParameterRequests(req)
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
