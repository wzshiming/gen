package swaggerui

import (
	"net/http"
	"os"

	"github.com/wzshiming/go-bindata/fs"
)

func HandleWith(f func(path string) ([]byte, error)) http.Handler {
	fun := f
	if fun == nil {
		fun = Asset
	} else {
		fun = func(path string) ([]byte, error) {
			data, err := f(path)
			if err == nil {
				return data, nil
			}
			return Asset(path)
		}
	}
	afs := &fs.AssetFS{
		Asset: fun,
		Index: "index.html",
	}
	return http.FileServer(afs)
}

func HandleWithFile(name string, data []byte) http.Handler {
	return HandleWith(func(path string) ([]byte, error) {
		if path == name {
			return data, nil
		}
		return nil, os.ErrNotExist
	})
}
