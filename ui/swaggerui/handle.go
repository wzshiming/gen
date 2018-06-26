package swaggerui

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"time"
)

var Handle = http.FileServer(SwaggerUI{})

type SwaggerUI struct{}

func (d SwaggerUI) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		return &Dir{name: name}, nil
	}
	data, err := Asset(name)
	if err != nil {
		return nil, os.ErrNotExist
	}

	return &File{
		name:   name,
		Reader: *bytes.NewReader(data),
	}, nil
}

type File struct {
	name string
	bytes.Reader
}

func (f *File) Close() error {
	return nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int64 {
	return f.Reader.Size()
}

func (f *File) Mode() os.FileMode {
	return 0600
}

func (f *File) ModTime() time.Time {
	return time.Time{}
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Sys() interface{} {
	return nil
}

type Dir struct {
	name string
}

func (f *Dir) Read(p []byte) (int, error) {
	return 0, nil
}

func (f *Dir) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *Dir) Close() error {
	return nil
}

func (f *Dir) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *Dir) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *Dir) Name() string {
	return f.name
}

func (f *Dir) Size() int64 {
	return 0
}

func (f *Dir) Mode() os.FileMode {
	return 0600
}

func (f *Dir) ModTime() time.Time {
	return time.Time{}
}

func (f *Dir) IsDir() bool {
	return true
}

func (f *Dir) Sys() interface{} {
	return nil
}
