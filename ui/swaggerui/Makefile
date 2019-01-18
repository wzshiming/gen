IMAGE := swaggerapi/swagger-ui
PWD   := $(shell pwd)

.PHONY: swaggerui.go swaggerui
swaggerui.go: swaggerui
	cd swaggerui && go-bindata --nomemcopy --ignore "\.map$$" --pkg swaggerui -o ../swaggerui.go ./

swaggerui:
	docker pull ${IMAGE}
	docker run --rm -v ${PWD}/swaggerui/:/swaggerui/ ${IMAGE} cp -r /usr/share/nginx/html/. /swaggerui
	sed -i '' -e 's#https://petstore.swagger.io/v2/swagger.json#openapi.json#' ./swaggerui/index.html

