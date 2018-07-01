package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) convertString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = %s`, out, in)
	return nil
}

func (g *GenRoute) convertPrtString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = &%s`, out, in)
	return nil
}

func (g *GenRoute) convertInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseInt(%s,0,0); err == nil {
	%s = %s(i)
}
`, in, out, name)
	return nil
}

func (g *GenRoute) convertPrtInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseInt(%s,0,0); err == nil {
	_i := %s(i)
	%s = &_i
}
`, in, name, out)
	return nil
}

func (g *GenRoute) convertUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseUint(%s,0,0); err == nil {
	%s = %s(i)
}
`, in, out, name)
	return nil
}

func (g *GenRoute) convertPrtUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	name := typ.Name
	g.buf.WriteFormat(`if i, err := strconv.ParseUint(%s,0,0); err == nil {
	_i := %s(i)
	%s = &_i
}
`, in, name, out)
	return nil
}

func (g *GenRoute) Convert(in, out string, typ *spec.Type) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	switch typ.Kind {
	case spec.Ptr:
		typ = typ.Elem
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		switch typ.Kind {
		case spec.String:
			return g.convertPrtString(in, out, typ)
		case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
			return g.convertPrtInt64(in, out, typ)
		case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
			return g.convertPrtUint64(in, out, typ)
		default:
		}
	default:
		switch typ.Kind {
		case spec.String:
			return g.convertString(in, out, typ)
		case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
			return g.convertInt64(in, out, typ)
		case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
			return g.convertUint64(in, out, typ)
		case spec.Slice:
			if typ.Elem.Kind == spec.Byte {
				g.buf.WriteFormat("%s := []byte(%s)\n", out, in)
				return nil
			}
		}
	}

	g.buf.WriteFormat("// Conversion of string to ")
	g.Types(typ)
	g.buf.WriteString(" is not supported.")
	g.buf.WriteFormat("\nvar %s ", out)
	g.Types(typ)

	return nil
}
