

SWAGGERUI_URL := https://codeload.github.com/swagger-api/swagger-ui/tar.gz/v3.20.1

swaggerui.go: swaggerui go-bindata
	cd swaggerui && go-bindata -nocompress -nomemcopy -pkg swaggerui -o ../swaggerui.go ./

swaggerui: swaggerui.tar.gz
	mkdir $@
	tar -xv --strip-components 2 -C $@ -f $? */dist/

swaggerui.tar.gz:
	curl -o $@ $(SWAGGERUI_URL)

go-bindata:
	go get github.com/wzshiming/go-bindata/cmd/go-bindata

