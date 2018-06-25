package client

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/wzshiming/gen"
	"github.com/wzshiming/namecase"
)

// GenClient is the generating generating
type GenClient struct {
	api *gen.API
	buf *bytes.Buffer
}

func NewGenClient(api *gen.API) *GenClient {
	return &GenClient{
		api: api,
		buf: bytes.NewBuffer(nil),
	}
}

func (g *GenClient) Generate() ([]byte, error) {
	g.buf.WriteString(`// Code generated; DO NOT EDIT.
package main

import (
	"github.com/wzshiming/requests"
)

`)
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
		g.buf.WriteString(commentLine(v.Description))
		g.buf.WriteString("type ")
		g.buf.WriteString(getName(k))
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
	operations := g.api.Operations
	for _, v := range operations {
		err = g.Operations(v)
		if err != nil {
			return err
		}
		err = g.FuncBody(v)
		if err != nil {
			return err
		}
		g.buf.WriteString("\n\n")
	}
	return
}

func (g *GenClient) FuncBody(oper *gen.Operation) (err error) {
	g.buf.WriteString("{\n")
	g.buf.WriteString("_resp, _err := NewRequests().\n")
	for _, v := range oper.Requests {
		req := v
		if req.Ref != "" {
			req = g.api.Requests[req.Ref]
		}

		switch req.In {
		case "header":
			g.buf.WriteString("SetHeader(\"")
			g.buf.WriteString(req.Name)
			g.buf.WriteString("\",fmt.Sprint(")
			g.buf.WriteString(req.Name)
			g.buf.WriteString(")).\n")
		case "cookie":
			// TODO
		case "path":
			g.buf.WriteString("SetPath(\"")
			g.buf.WriteString(req.Name)
			g.buf.WriteString("\",fmt.Sprint(")
			g.buf.WriteString(req.Name)
			g.buf.WriteString(")).\n")
		case "query":
			g.buf.WriteString("SetQuery(\"")
			g.buf.WriteString(req.Name)
			g.buf.WriteString("\",fmt.Sprint(")
			g.buf.WriteString(req.Name)
			g.buf.WriteString(")).\n")
		case "body":
			switch req.Content {
			case "json":
				g.buf.WriteString("SetJSON(")
				g.buf.WriteString(req.Name)
				g.buf.WriteString(").\n")
			case "xml":
				g.buf.WriteString("SetJSON(")
				g.buf.WriteString(req.Name)
				g.buf.WriteString(").\n")
			}
		}
	}
	g.buf.WriteString(namecase.ToPascal(oper.Method) + "(\"" + oper.Path + "\")\n")

	g.buf.WriteString(`
if _err != nil {
	err = _err
	return
}
`)

	g.buf.WriteString("switch _resp.StatusCode() {\n")

	for _, v := range oper.Responses {
		resp := v
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		g.buf.WriteString("case ")
		g.buf.WriteString(resp.Code)
		g.buf.WriteString(":\n")
		// TODO
	}

	g.buf.WriteString(`default:
	err = fmt.Errorf("Undefined code %d:", _resp.StatusCode())
}
return
}
`)
	return nil
}

func (g *GenClient) Operations(oper *gen.Operation) (err error) {

	g.buf.WriteString(commentLine(oper.Description))
	g.buf.WriteString("func ")

	if oper.Type != nil {
		g.buf.WriteByte('(')
		err = g.Types(oper.Type)
		if err != nil {
			return err
		}
		g.buf.WriteByte(')')
	}
	g.buf.WriteString(getName(oper.Name))
	g.buf.WriteByte('(')
	for i, v := range oper.Requests {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		err = g.Requests(v)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(")(")
	needErr := false
	for i, v := range oper.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		err = g.Responses(v)
		if err != nil {
			return err
		}
	}
	if len(oper.Responses) == 0 {
		needErr = true
	} else {
		resp := oper.Responses[len(oper.Responses)-1]
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		if resp.Type.Type != "error" {
			needErr = true
		}
	}
	if needErr {
		if len(oper.Responses) != 0 {
			g.buf.WriteByte(',')
		}
		g.buf.WriteString("err error")
	}
	g.buf.WriteByte(')')
	return
}

func (g *GenClient) Requests(req *gen.Request) (err error) {
	if req.Ref != "" {
		return g.Requests(g.api.Requests[req.Ref])
	}
	g.buf.WriteString(req.Name)
	g.buf.WriteByte(' ')
	err = g.Types(req.Type)
	if err != nil {
		return err
	}
	if req.Description != "" {
		g.buf.WriteString("/* ")
		g.buf.WriteString(mergeLine(req.Description))
		g.buf.WriteString(" */")
	}
	return nil
}

func (g *GenClient) Responses(req *gen.Response) (err error) {
	if req.Ref != "" {
		return g.Responses(g.api.Responses[req.Ref])
	}
	g.buf.WriteString(req.Name)
	g.buf.WriteByte(' ')
	err = g.Types(req.Type)
	if err != nil {
		return err
	}
	if req.Description != "" {
		g.buf.WriteString("/* ")
		g.buf.WriteString(mergeLine(req.Description))
		g.buf.WriteString(" */")
	}
	return nil
}

func (g *GenClient) Types(typ *gen.Type) (err error) {
	if typ.Ref != "" {
		g.buf.WriteString(getName(typ.Ref))
		return nil
	}
	switch typ.Type {
	case "ptr":
		g.buf.WriteByte('*')
		err := g.Types(typ.Elem)
		if err != nil {
			return err
		}
	case "slice":
		g.buf.WriteString("[]")
		err := g.Types(typ.Elem)
		if err != nil {
			return err
		}
	case "array":
		g.buf.WriteByte('[')
		g.buf.WriteString(strconv.Itoa(typ.Len))
		g.buf.WriteByte(']')
		err := g.Types(typ.Elem)
		if err != nil {
			return err
		}
	case "map":
		g.buf.WriteString("map[")
		err := g.Types(typ.Key)
		if err != nil {
			return err
		}
		g.buf.WriteByte(']')
		err = g.Types(typ.Elem)
		if err != nil {
			return err
		}
	case "struct":
		g.buf.WriteString("struct {")
		if len(typ.Fields) != 0 {
			g.buf.WriteByte('\n')
		}
		for _, v := range typ.Fields {
			if !v.Anonymous {
				g.buf.WriteString(v.Name)
				g.buf.WriteByte(' ')
			}
			err := g.Types(v.Type)
			if err != nil {
				return err
			}
			if v.Tag != "" {
				g.buf.WriteByte(' ')
				g.buf.WriteByte('`')
				g.buf.WriteString(string(v.Tag))
				g.buf.WriteByte('`')
			}
			g.buf.WriteString("// ")
			g.buf.WriteString(mergeLine(v.Description))
			g.buf.WriteByte('\n')
		}
		g.buf.WriteByte('}')
	default:
		g.buf.WriteString(typ.Type)
	}
	return
}

func getName(name string) string {
	i := strings.Index(name, ".")
	if i == -1 {
		return name
	}
	return name[:i]
}

func mergeLine(t string) string {
	return strings.TrimSpace(strings.Replace(t, "\n", " ", -1))
}

func commentLine(t string) string {
	return "// " + strings.Join(strings.Split(strings.TrimSpace(t), "\n"), "\n// ") + "\n"
}
