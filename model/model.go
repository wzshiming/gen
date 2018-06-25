package model

import (
	"bytes"
	"sort"
	"strings"

	"github.com/wzshiming/gen"
)

// GenModel is the parse type generating definitions
type GenModel struct {
	api *gen.API
	buf *bytes.Buffer
}

func NewGenModel(api *gen.API) *GenModel {
	return &GenModel{
		api: api,
		buf: bytes.NewBuffer(nil),
	}
}

func (g *GenModel) Generate() ([]byte, error) {
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

func (g *GenModel) GenerateSchemas() (err error) {
	schemas := g.api.Types
	ks := make([]string, 0, len(schemas))
	for k, _ := range schemas {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		v := schemas[k]

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

func (g *GenModel) GenerateOperations() (err error) {
	operations := g.api.Operations
	for _, v := range operations {
		err = g.Operations(v)
		if err != nil {
			return err
		}
		g.buf.WriteString("\n\n")
	}
	return
}

func (g *GenModel) Operations(oper *gen.Operation) (err error) {

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
			g.buf.WriteString(", ")
		}
		err = g.Requests(v)
		if err != nil {
			return err
		}
	}
	g.buf.WriteString(") (")
	for i, v := range oper.Responses {
		if i != 0 {
			g.buf.WriteString(", ")
		}
		err = g.Responses(v)
		if err != nil {
			return err
		}
	}
	g.buf.WriteByte(')')
	g.buf.WriteString("{}")
	return
}

func (g *GenModel) Requests(req *gen.Request) (err error) {
	if req.Ref != "" {
		return g.Requests(g.api.Requests[req.Ref])
	}
	g.buf.WriteString(req.Name)
	g.buf.WriteByte(' ')
	return g.Types(req.Type)
}

func (g *GenModel) Responses(req *gen.Response) (err error) {
	if req.Ref != "" {
		return g.Responses(g.api.Responses[req.Ref])
	}
	g.buf.WriteString(req.Name)
	g.buf.WriteByte(' ')
	return g.Types(req.Type)
}

func (g *GenModel) Types(typ *gen.Type) (err error) {
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
	case "array":
		g.buf.WriteString("[]")
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
			g.buf.WriteString(v.Name)
			g.buf.WriteByte(' ')
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
