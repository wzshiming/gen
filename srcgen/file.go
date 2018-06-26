package srcgen

import (
	"go/format"
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
		buf.WriteFormat(")")
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
