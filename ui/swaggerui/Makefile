IMAGE := swaggerapi/swagger-ui
PWD   := $(shell pwd)

.PHONY: swaggerui.go swaggerui
swaggerui.go: swaggerui go-bindata
	cd swaggerui && go-bindata --nomemcopy --pkg swaggerui -o ../swaggerui.go ./

swaggerui:
	docker pull ${IMAGE}
	docker run -v ${PWD}/swaggerui/:/swaggerui/ ${IMAGE} cp -r /usr/share/nginx/html/. /swaggerui
	rm -f ./swaggerui/*.map
	sed -i '' -e 's#https://petstore.swagger.io/v2/swagger.json#openapi.json#' ./swaggerui/index.html

go-bindata:
	go get github.com/wzshiming/go-bindata/cmd/go-bindata

