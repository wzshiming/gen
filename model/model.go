package model

import (
	"strconv"
	"strings"

	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gen/srcgen"
	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/namecase"
)

// GenModel is the generating generating
type GenModel struct {
	pkgpath string
	api     *spec.API
	buf     *srcgen.File
}

func NewGenModel(api *spec.API, buf *srcgen.File, pkgpath string) *GenModel {
	return &GenModel{
		api:     api,
		buf:     buf,
		pkgpath: pkgpath,
	}
}

func (g *GenModel) TypesZero(typ *spec.Type) (err error) {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	switch typ.Kind {
	case spec.Ptr, spec.Slice, spec.Array, spec.Map, spec.Error, spec.Chan, spec.Interface:
		g.buf.WriteString("nil")
	case spec.String:
		g.buf.WriteString("\"\"")
	case spec.Struct:
		g.Types(typ)
		g.buf.WriteString("{}")
	case spec.Int, spec.Int8, spec.Int16, spec.Int32, spec.Int64,
		spec.Uint, spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64,
		spec.Float32, spec.Float64, spec.Complex64, spec.Complex128,
		spec.Byte, spec.Rune:
		g.buf.WriteString("0")
	case spec.Bool:
		g.buf.WriteString("false")
	default:
		g.buf.WriteString(strings.ToLower(typ.Kind.String()))
	}
	return nil
}

func (g *GenModel) PkgPath(path string) bool {
	if g.pkgpath == "" || path == g.pkgpath {
		return false
	}
	pkgname := namecase.ToCamel(path)
	g.buf.AddImport(pkgname, path)
	g.buf.WriteFormat("%s.", pkgname)
	return true
}

func (g *GenModel) Paths(typ *spec.Type) bool {
	reftyp, ok := g.api.Types[typ.Ref]
	if !ok {
		return false
	}
	return g.PkgPath(reftyp.PkgPath)
}

func (g *GenModel) Types(typ *spec.Type) (err error) {
	if typ.Ref != "" {
		g.Paths(typ)
		g.buf.WriteString(utils.GetName(typ.Ref))
		return nil
	}
	switch typ.Kind {
	case spec.Ptr:
		g.buf.WriteByte('*')
		err := g.Types(typ.Elem)
		if err != nil {
			return err
		}
	case spec.Slice:
		g.buf.WriteString("[]")
		err := g.Types(typ.Elem)
		if err != nil {
			return err
		}
	case spec.Array:
		g.buf.WriteByte('[')
		g.buf.WriteString(strconv.Itoa(typ.Len))
		g.buf.WriteByte(']')
		err := g.Types(typ.Elem)
		if err != nil {
			return err
		}
	case spec.Map:
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
	case spec.Struct:
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
		g.buf.WriteString(strings.ToLower(typ.Kind.String()))
	}
	return
}
