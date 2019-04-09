package model

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenModel) convertToString(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s = `, out)
	g.Types(typ)
	g.buf.WriteFormat(`(%s)`, in)
	return nil
}

func (g *GenModel) convertToInt64(in, out string, typ *spec.Type) error {
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

func (g *GenModel) convertToUint64(in, out string, typ *spec.Type) error {
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

func (g *GenModel) convertToBytes(in, out string, typ *spec.Type) error {
	g.buf.WriteFormat(`%s := `, out)
	g.Types(typ)
	g.buf.WriteFormat(`(%s)
`, in)
	return nil
}

func (g *GenModel) convertToSlice(in, out string, typ *spec.Type) error {
	g.buf.AddImport("", "strconv")
	g.buf.AddImport("", "strings")
	g.buf.WriteFormat(`
	_list_%s := strings.Split(%s, ",")
`, in, in)
	g.buf.WriteFormat(`%s = make([]`, out)
	g.Types(typ)
	g.buf.WriteFormat(`, 0, len(_list_%s))
`, in)

	g.buf.WriteFormat(`
	for _, _%s := range _list_%s {
		var _%s `, in, in, out)
	g.Types(typ)
	g.buf.WriteFormat(`
`)
	err := g.convertTo("_"+in, "_"+out, typ)
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

func (g *GenModel) convertTo(in, out string, typ *spec.Type) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}

	if typ.Attr.Has(spec.AttrTextUnmarshaler) {
		g.buf.AddImport("", "unsafe")
		g.buf.AddImport("", "net/http")
		g.buf.WriteFormat(`
	if %s != "" {
		err = %s.UnmarshalText(*(*[]byte)(unsafe.Pointer(&%s)))
	}
`, in, out, in)
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
		return g.convertToString(in, out, typ)
	case spec.Int8, spec.Int16, spec.Int32, spec.Int64, spec.Int:
		return g.convertToInt64(in, out, typ)
	case spec.Uint8, spec.Uint16, spec.Uint32, spec.Uint64, spec.Uint:
		return g.convertToUint64(in, out, typ)
	case spec.Slice:
		switch typ.Elem.Kind {
		case spec.Byte, spec.Rune:
			return g.convertToBytes(in, out, typ)
		}
		return g.convertToSlice(in, out, typ.Elem)
	case spec.Array:
		return g.convertToSlice(in, out, typ.Elem)
	}

	g.buf.WriteFormat("// Conversion of string to ")
	g.Types(typ)
	g.buf.WriteString(" is not supported.")

	return nil
}

func (g *GenModel) ConvertTo(in, out string, typ *spec.Type) error {
	return g.convertTo(in, out, typ)
}

func (g *GenModel) convertToMulti(in, out string, typ *spec.Type) error {
	if typ.Ref != "" {
		typ = g.api.Types[typ.Ref]
	}
	g.buf.WriteFormat(`
	if len(%s) == 0 {
		return
	}
`, in)

	switch typ.Kind {
	default:
		return g.convertTo(in+"[0]", out, typ)
	case spec.Slice, spec.Array:
	}

	g.buf.WriteFormat(`%s = make(`, out)
	g.Types(typ)
	g.buf.WriteFormat(`, 0, len(%s))
`, in)

	g.buf.WriteFormat(`
	for _, _%s := range %s {
		var _%s `, in, in, out)
	g.Types(typ.Elem)
	g.buf.WriteFormat(`
`)
	err := g.convertTo("_"+in, "_"+out, typ.Elem)
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

func (g *GenModel) ConvertToMulti(in, out string, typ *spec.Type) error {
	return g.convertToMulti(in, out, typ)
}
