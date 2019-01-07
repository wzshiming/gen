package crud

import (
	"bytes"
	"strings"
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

const suf = ".go.tpl"

func TplNames() []string {
	an := tpl.AssetNames()
	names := make([]string, 0, len(an))
	for _, name := range an {
		if strings.HasSuffix(name, suf) {
			names = append(names, name[:len(name)-len(suf)])
		}
	}
	return names
}

func (g *GenCrud) Generate(tplname, pkgname, typname string) ([]byte, error) {
	file, err := tpl.Asset(tplname + suf)
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
