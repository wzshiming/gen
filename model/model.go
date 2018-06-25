package model

import (
	"bytes"
	"strconv"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/utils"
)

// GenModel is the generating generating
type GenModel struct {
	api *spec.API
	buf *bytes.Buffer
}

func NewGenModel(api *spec.API, buf *bytes.Buffer) *GenModel {
	return &GenModel{
		api: api,
		buf: buf,
	}
}

func (g *GenModel) Types(typ *spec.Type) (err error) {
	if typ.Ref != "" {
		g.buf.WriteString(utils.GetName(typ.Ref))
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
			g.buf.WriteString(utils.MergeLine(v.Description))
			g.buf.WriteByte('\n')
		}
		g.buf.WriteByte('}')
	default:
		g.buf.WriteString(typ.Type)
	}
	return
}
