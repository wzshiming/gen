package crud

import (
	"bytes"
	"text/template"
	"unsafe"

	"github.com/wzshiming/gen/crud/tpl"
	"github.com/wzshiming/namecase"
)

type GenCrud struct {
}

func NewGenCrud() *GenCrud {
	return &GenCrud{}
}

func (g *GenCrud) Generate(tplname, pkgname, typname string) ([]byte, error) {
	file, err := tpl.Asset(tplname + ".go.tpl")
	if err != nil {
		return nil, err
	}

	fb := *(*string)(unsafe.Pointer(&file))
	temp, err := template.New("").Delims("<", ">").Parse(fb)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	err = temp.Execute(buf, map[string]string{
		"Package":    pkgname,
		"Original":   typname,
		"UpperHump":  namecase.ToUpperHump(typname),
		"LowerHump":  namecase.ToLowerHump(typname),
		"LowerSnake": namecase.ToLowerSnake(typname),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
