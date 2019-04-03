package client

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenClient) convertString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = `, out)
	g.buf.WriteFormat(`string(%s)
`, in)
	return nil
}

func (g *GenClient) convertFromInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	%s = `, out)
	g.buf.WriteFormat(`strconv.FormatInt(int64(%s), 10)
`, in)

	return nil
}

func (g *GenClient) convertFromUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	%s = `, out)
	g.buf.WriteFormat(`strconv.FormatUint(uint64(%s), 10)
`, in)
	return nil
}

func (g *GenClient) convertFromBytes(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = `, out)
	g.buf.WriteFormat(`%s
`, in)
	return nil
}

func (g *GenClient) convertFromSlice(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strings")

	g.buf.WriteFormat(`_list_%s := make([]string`, out)
	g.buf.WriteFormat(`, 0, len(%s))
`, in)

	g.buf.WriteFormat(`
	for _, _%s := range %s {
		var _%s string
`, in, in, out)

	err := g.convertFrom("_"+in, "_"+out, typ)
	if err != nil {
		return err
	}
	g.buf.WriteFormat(`
		_list_%s = append(_list_%s, _%s)
	}

	%s = strings.Join(_list_%s, ",")
`, out, out, out, out, out)

	return nil
}

func (g *GenClient) convertFrom(in, out string, typ *spec.Type) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	if typ.Attr.Has(spec.AttrTextMarshaler) || typ.Kind == spec.Time {
		g.buf.AddImport("", "unsafe")
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	var _%s []byte
	_%s, err = %s.MarshalText()
	if err != nil {
		return
	}
	%s = *(*string)(unsafe.Pointer(&_%s))

`, out, out, in, out, out)
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
		return g.convertFromInt64(in, out, typ)
	case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
		return g.convertFromUint64(in, out, typ)
	case spec.Slice:
		switch typ.Elem.Kind {
		case spec.Byte, spec.Rune:
			return g.convertFromBytes(in, out, typ)
		}
		return g.convertFromSlice(in, out, typ.Elem)
	case spec.Array:
		return g.convertFromSlice(in, out, typ.Elem)
	}

	g.buf.WriteFormat("// Conversion of string from ")
	g.Types(typ)
	g.buf.WriteString(" is not supported.")

	return nil
}
