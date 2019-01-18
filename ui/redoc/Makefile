IMAGE := redocly/redoc
PWD   := $(shell pwd)

.PHONY: redoc.go redoc
redoc.go: redoc
	cd redoc && go-bindata --nomemcopy --ignore "\.map$$" --pkg redoc -o ../redoc.go ./

redoc:
	docker pull ${IMAGE}
	docker run --rm -v ${PWD}/redoc/:/redoc/ ${IMAGE} cp -r /usr/share/nginx/html/. /redoc
	sed -i '' -e 's|%PAGE_TITLE%|ReDoc|g' ./redoc/index.html
	sed -i '' -e 's|%PAGE_FAVICON%|favicon.png|g' ./redoc/index.html
	sed -i '' -e 's|%SPEC_URL%|./openapi.json|g' ./redoc/index.html
	sed -i '' -e 's|%REDOC_OPTIONS%||g' ./redoc/index.html
	sed -i '' -e 's|/redoc.standalone.js|./redoc.standalone.js|g' ./redoc/index.html

