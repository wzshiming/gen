package model

import (
	"path"
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
	typ0 := typ
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	switch typ.Kind {
	case spec.Ptr, spec.Slice, spec.Map, spec.Error, spec.Chan, spec.Interface:
		g.buf.WriteString("nil")
	case spec.Array:
		g.buf.WriteString("(")
		err = g.Types(typ0)
		if err != nil {
			return err
		}
		g.buf.WriteString("{})")
	case spec.String:
		g.buf.WriteString("\"\"")
	case spec.Struct:
		err = g.Types(typ0)
		if err != nil {
			return err
		}
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
	const (
		vendor = "/vendor/"
	)
	if i := strings.LastIndex(path, vendor); i != -1 {
		path = path[i+len(vendor):]
	}
	if g.pkgpath == "" || path == g.pkgpath {
		return false
	}
	pkgname := namecase.ToCamel(path)
	g.buf.AddImport(pkgname, path)
	g.buf.WriteFormat("%s.", pkgname)
	return true
}

func (g *GenModel) Paths(typ *spec.Type) bool {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	return g.PkgPath(typ.PkgPath)
}

func (g *GenModel) Ptr(typ *spec.Type) bool {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	return typ.Kind != spec.Interface
}

func (g *GenModel) PtrTypes(typ *spec.Type) (err error) {
	if typ.Ref != "" {
		if g.Ptr(typ) {
			g.buf.WriteString("*")
		}
		g.Paths(typ)
		g.buf.WriteString(getName(typ.Ref))
		return nil
	}

	if typ.Name != "" && strings.ToLower(typ.Kind.String()) != typ.Name {
		g.Paths(typ)
		g.buf.WriteString(typ.Name)
		return nil
	}

	if g.Ptr(typ) {
		g.buf.WriteString("*")
	}
	return g.TypesDefine(typ)
}

func (g *GenModel) Types(typ *spec.Type) (err error) {

	if g.typeRoot(typ) {
		return
	}

	if typ.Ref != "" {
		g.Paths(typ)
		g.buf.WriteString(getName(typ.Ref))
		return nil
	}

	if typ.Name != "" && strings.ToLower(typ.Kind.String()) != typ.Name {
		g.Paths(typ)
		g.buf.WriteString(typ.Name)
		return nil
	}

	return g.TypesDefine(typ)
}

func (g *GenModel) typeRoot(typ *spec.Type) bool {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	if typ.Attr.Has(spec.AttrRoot) {
		g.buf.AddImport("", typ.PkgPath)
		_, pkgname := path.Split(typ.PkgPath)
		g.buf.WriteFormat("%s.%s", pkgname, typ.Name)
		return true
	}
	return false
}

func (g *GenModel) TypesDefine(typ *spec.Type) (err error) {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	if g.typeRoot(typ) {
		return
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
	case spec.Interface:
		g.buf.WriteString("interface{}")
	default:
		g.buf.WriteString(strings.ToLower(typ.Kind.String()))
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
