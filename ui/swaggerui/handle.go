package swaggerui

import (
	"net/http"

	"github.com/wzshiming/go-bindata/fs"
)

var (
	afs = &fs.AssetFS{
		Asset: Asset,
		Index: "index.html",
	}
	Handle = http.FileServer(afs)
)
