package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) convertString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = `, out)
	g.Types(typ)
	g.buf.WriteFormat(`(%s)`, in)
	return nil
}

func (g *GenRoute) convertPrtString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`__%s := `, out)
	g.Types(typ)
	g.buf.WriteFormat(`(%s)
	%s = &__%s
`, in, out, out)
	return nil
}

func (g *GenRoute) convertInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	if _%s, err := strconv.ParseInt(%s, 0, 0); err == nil {
		%s = `, out, in, out)
	g.Types(typ)
	g.buf.WriteFormat(`(_%s)
	}
`, out)
	return nil
}

func (g *GenRoute) convertPrtInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	if _%s, err := strconv.ParseInt(%s, 0, 0); err == nil {
		__%s := `, out, in, out)
	g.Types(typ)
	g.buf.WriteFormat(`(%s)
		%s = &__%s
	}
`, out, out, out)
	return nil
}

func (g *GenRoute) convertUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	if _%s, err := strconv.ParseUint(%s, 0, 0); err == nil {
		%s = `, out, in, out)
	g.Types(typ)
	g.buf.WriteFormat(`(_%s)
	}
`, out)
	return nil
}

func (g *GenRoute) convertPrtUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	if _%s, err := strconv.ParseUint(%s, 0, 0); err == nil {
		__%s := `, out, in, out)
	g.Types(typ)
	g.buf.WriteFormat(`(_%s)
		%s = &__%s
	}
`, out, out, out)
	return nil
}

func (g *GenRoute) convertBytes(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s := `, out)
	g.Types(typ)
	g.buf.WriteFormat(`(%s)
`, in)
	return nil
}

func (g *GenRoute) Convert(in, out string, typ *spec.Type) error {
	ptyp := typ
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	if typ.IsTextUnmarshaler || typ.Kind == spec.Time {
		g.buf.AddImport("", "unsafe")
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	err = %s.UnmarshalText(*(*[]byte)(unsafe.Pointer(&%s)))
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
`, out, in)
		return nil
	}

	switch typ.Kind {
	case spec.Ptr:
		typ = typ.Elem
		ptyp = typ
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		switch typ.Kind {
		case spec.String:
			return g.convertPrtString(in, out, ptyp)
		case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
			return g.convertPrtInt64(in, out, ptyp)
		case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
			return g.convertPrtUint64(in, out, ptyp)
		default:
		}
	default:
		switch typ.Kind {
		case spec.String:
			return g.convertString(in, out, ptyp)
		case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
			return g.convertInt64(in, out, ptyp)
		case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
			return g.convertUint64(in, out, ptyp)
		case spec.Slice:
			switch typ.Elem.Kind {
			case spec.Byte, spec.Rune:
				return g.convertBytes(in, out, ptyp)
			}
		}
	}

	g.buf.WriteFormat("// Conversion of string to ")
	g.Types(typ)
	g.buf.WriteString(" is not supported.")

	return nil
}
