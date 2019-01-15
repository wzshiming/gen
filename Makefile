
test:
	go test ./...

generate: go-bindata
	make -C crud/tpl
	make -C ui/swaggerui
	make -C ui/redoc
	make -C examples

go-bindata:
	go get github.com/wzshiming/go-bindata/cmd/go-bindata

