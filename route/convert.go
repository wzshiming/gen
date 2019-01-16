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

func (g *GenRoute) convertBytes(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s := `, out)
	g.Types(typ)
	g.buf.WriteFormat(`(%s)
`, in)
	return nil
}

func (g *GenRoute) Convert(in, out string, typ *spec.Type) error {
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

	if typ.Name == "" && typ.Kind == spec.Ptr {
		out0 := out
		out = "_" + out
		typ = typ.Elem
		g.buf.WriteFormat(`
	var %s `, out)
		g.Types(typ)
		defer func() {
			g.buf.WriteFormat(`
	%s = &%s
`, out0, out)
		}()
	}

	switch typ.Kind {
	case spec.String:
		return g.convertString(in, out, typ)
	case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
		return g.convertInt64(in, out, typ)
	case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
		return g.convertUint64(in, out, typ)
	case spec.Slice:
		switch typ.Elem.Kind {
		case spec.Byte, spec.Rune:
			return g.convertBytes(in, out, typ)
		}
	}

	g.buf.WriteFormat("// Conversion of string to ")
	g.Types(typ)
	g.buf.WriteString(" is not supported.")

	return nil
}
