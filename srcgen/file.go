package srcgen

import (
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/wzshiming/namecase"
	"golang.org/x/tools/imports"
)

type File struct {
	srcgen
	filename    string
	packname    string
	imports     map[string]string
	buildIgnore bool
}

func NewFile() *File {
	return &File{}
}

func (f *File) Save() error {
	if f.filename == "" {
		f.filename = "auto_gen.go"
	}

	if f.packname == "" {
		dir := filepath.Dir(f.filename)
		b, err := build.ImportDir(dir, 0)
		if err != nil {
			_, dir = filepath.Split(dir)
			if dir == "." || dir == ".." || dir == "" {
				f.packname = "main"
			} else {
				f.packname = dir
			}
		} else {
			f.packname = b.Name
		}
		f.packname = namecase.ToLowerSnake(f.packname)
	}

	err := os.MkdirAll(filepath.Dir(f.filename), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.filename, f.Bytes(), 0666)
}

func (f *File) WithPackname(packname string) *File {
	f.packname = packname
	return f
}

func (f *File) WithFilename(filename string) *File {
	f.filename = filename
	return f
}

func (f *File) AddImport(aliase, importpath string) *File {
	if f.imports == nil {
		f.imports = map[string]string{}
	}
	if aliase == "" {
		_, aliase = path.Split(importpath)
		aliase = strings.SplitN(aliase, ".", 2)[0]
	}
	if aliase == "_" {
		_, ok := f.imports[importpath]
		if ok {
			return f
		}
	}
	f.imports[importpath] = aliase
	return f
}

func (f *File) SetBuildIgnore(b bool) *File {
	f.buildIgnore = b
	return f
}

func (f *File) Bytes() []byte {
	buf := srcgen{}
	if f.buildIgnore {
		buf.WriteFormat(`// +build ignore
`)
	}
	buf.WriteFormat(`// Code generated; DO NOT EDIT.
// file %s

package %s

`, f.filename, f.packname)

	if len(f.imports) != 0 {
		buf.WriteFormat("import(\n")
		for path, aliase := range f.imports {
			buf.WriteFormat("%s \"%s\"\n", aliase, path)
		}
		buf.WriteFormat(")\n\n")
	}

	buf.WriteString(f.srcgen.String())

	data := buf.Bytes()

	dataf, err := imports.Process(f.filename, data, nil)
	if err != nil {
		return data
	}
	return dataf
}

func (f *File) String() string {
	data := f.Bytes()
	return *(*string)(unsafe.Pointer(&data))
}
