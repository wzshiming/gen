

SWAGGERUI_URL := https://codeload.github.com/swagger-api/swagger-ui/tar.gz/v3.20.1

swaggerui.go: swaggerui go-bindata
	cd swaggerui && sed -i '' -e 's#https://petstore.swagger.io/v2/swagger.json#openapi.json#' index.html && go-bindata --nomemcopy --pkg swaggerui -o ../swaggerui.go ./

swaggerui: swaggerui.tar.gz
	mkdir $@
	tar -xv --strip-components 2 -C $@ -f $? */dist/

swaggerui.tar.gz:
	curl -o $@ $(SWAGGERUI_URL)

go-bindata:
	go get github.com/wzshiming/go-bindata/cmd/go-bindata

