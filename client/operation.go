package client

import (
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
)

func (g *GenClient) generateClient() (err error) {
	g.buf.AddImport("", "github.com/wzshiming/requests")
	g.buf.WriteFormat(`
	var Client = requests.NewClient().NewRequest()
`)

	operations := g.api.Operations
	for _, v := range operations {
		err = g.generateOperations(v)
		if err != nil {
			return err
		}
		err = g.generateFuncBody(v)
		if err != nil {
			return err
		}
		g.buf.WriteString("\n\n")
	}
	return
}

func (g *GenClient) mergeMiddlewareRequests(sreq []*spec.Request) (reqs []*spec.Request, err error) {
	for _, req := range sreq {
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}
		switch req.In {
		case "security":
			// No action
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
					r, err := g.mergeMiddlewareRequests(v.Requests)
					if err != nil {
						return nil, err
					}
					reqs = append(reqs, r...)
				}
			}

		case "header", "path", "query", "body":
			reqs = append(reqs, req)
		}
	}
	return reqs, nil
}

func (g *GenClient) generateOperations(oper *spec.Operation) (err error) {

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
	g.buf.WriteString(g.getFuncName(oper))
	g.buf.WriteByte('(')
	reqs, err := g.mergeMiddlewareRequests(oper.Requests)
	if err != nil {
		return err
	}
	for i, req := range reqs {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		typ := ""
		switch req.Content {
		case "file", "image":
			g.buf.AddImport("", "io")
			typ = "io.Reader"
		}
		err = g.generateParameterRequests(req, typ)
		if err != nil {
			return err
		}

	}
	g.buf.WriteString(")(")

	for i, v := range oper.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		err = g.generateParameterResponses(v)
		if err != nil {
			return err
		}
	}

	g.buf.WriteByte(')')
	return
}
