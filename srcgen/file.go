package srcgen

import (
	"go/build"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"unsafe"
)

type File struct {
	srcgen
	filename string
	packname string
	imports  map[string]string
}

func NewFile() *File {
	return &File{}
}

func (f *File) Save() error {
	if f.filename == "" {
		f.filename = "auto_gen.go"
	}

	if f.packname == "" {
		b, err := build.ImportDir(filepath.Dir(f.filename), 0)
		if err != nil {
			return err
		}
		f.packname = b.Name
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

func (f *File) AddImport(aliase, path string) *File {
	if f.imports == nil {
		f.imports = map[string]string{}
	}
	f.imports[path] = aliase
	return f
}

func (f *File) Bytes() []byte {
	buf := srcgen{}
	buf.WriteFormat(`// Code generated; DO NOT EDIT.

// %s
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

	dataf, err := format.Source(data)
	if err != nil {
		return data
	}
	return dataf
}

func (f *File) String() string {
	data := f.Bytes()
	return *(*string)(unsafe.Pointer(&data))
}
