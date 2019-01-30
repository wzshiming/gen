
test:
	go test ./...

generate: go-bindata
	make -C crud/tpl
	make -C examples

go-bindata:
	go get github.com/wzshiming/go-bindata/cmd/go-bindata

