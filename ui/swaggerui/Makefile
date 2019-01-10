
swaggerui.go: swaggerui go-bindata
	cd swaggerui && go-bindata --nomemcopy --pkg swaggerui -o ../swaggerui.go ./

swaggerui:
	docker run -v $(PWD)/swaggerui/:/swaggerui/ swaggerapi/swagger-ui cp -r /usr/share/nginx/html/. /swaggerui
	rm -f ./swaggerui/*.map
	sed -i '' -e 's#https://petstore.swagger.io/v2/swagger.json#openapi.json#' ./swaggerui/index.html

go-bindata:
	go get github.com/wzshiming/go-bindata/cmd/go-bindata

