

SWAGGERUI_URL := https://codeload.github.com/swagger-api/swagger-ui/tar.gz/v3.17.1
BINDATA       := $(GOPATH)/bin/go-bindata

swaggerui.go: swaggerui $(BINDATA)
	cd swaggerui && $(BINDATA) -pkg swaggerui -o ../swaggerui.go ./

swaggerui: swaggerui.tar.gz
	mkdir $@
	tar -xv --wildcards --strip-components 2 -C $@ -f $? */dist/

swaggerui.tar.gz:
	curl -o $@ $(SWAGGERUI_URL)

$(BINDATA):
	go get github.com/go-bindata/go-bindata/go-bindata