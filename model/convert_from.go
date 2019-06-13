package model

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenModel) convertFromString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = `, out)
	g.buf.WriteFormat(`string(%s)
`, in)
	return nil
}

func (g *GenModel) convertFromInt64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	%s = `, out)
	g.buf.WriteFormat(`strconv.FormatInt(int64(%s), 10)
`, in)

	return nil
}

func (g *GenModel) convertFromUint64(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.WriteFormat(`
	%s = `, out)
	g.buf.WriteFormat(`strconv.FormatUint(uint64(%s), 10)
`, in)
	return nil
}

func (g *GenModel) convertFromBytes(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = `, out)
	g.buf.WriteFormat(`%s
`, in)
	return nil
}

func (g *GenModel) convertFromSlice(in, out string, typ *spec.Type) error {
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

func (g *GenModel) convertFrom(in, out string, typ *spec.Type) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	if typ.Attr.Has(spec.AttrTextMarshaler) {
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
		return g.convertFromString(in, out, typ)
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

func (g *GenModel) ConvertFrom(in, out string, typ *spec.Type) error {
	return g.convertFrom(in, out, typ)
}

func (g *GenModel) convertFromMulti(in, out string, typ *spec.Type, explode bool) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	g.buf.WriteFormat(`
	if len(%s) == 0 {
		return
	}
`, in)

	if !explode || (typ.Kind != spec.Slice && typ.Kind != spec.Array) {
		g.buf.WriteFormat(`
		_%s_0 := %s[0]
		`, in, in)
		return g.convertFrom("_"+in+"_0", out, typ)
	}

	g.buf.WriteFormat(`%s = make([]string, 0, len(%s))
`, out, in)

	g.buf.WriteFormat(`
	for _, _%s := range %s {
		var _%s string
	`, in, in, out)

	err := g.convertFrom("_"+in, "_"+out, typ.Elem)
	if err != nil {
		return err
	}
	g.buf.WriteFormat(`
		if err != nil {
			break
		}
		%s = append(%s, _%s)
	}
`, out, out, out)

	return nil
}

func (g *GenModel) ConvertFromMulti(in, out string, typ *spec.Type, explode bool) error {
	return g.convertFromMulti(in, out, typ, explode)
}
