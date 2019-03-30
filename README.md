# Gen - Tools for generating source code for microservices

Just write normal functions, and Gen generates efficient routing source code and documentation for it
Because the source code is generated, none of this affects runtime performance  
The differences caused by each change in the tool are shown directly in the generated source code  
generating clients is also supported  

[![Build Status](https://travis-ci.org/wzshiming/gen.svg?branch=master)](https://travis-ci.org/wzshiming/gen)
[![Go Report Card](https://goreportcard.com/badge/github.com/wzshiming/gen)](https://goreportcard.com/report/github.com/wzshiming/gen)
[![GitHub license](https://img.shields.io/github/license/wzshiming/gen.svg)](https://github.com/wzshiming/gen/blob/master/LICENSE)

- [English](https://github.com/wzshiming/gen/blob/master/README.md)
- [简体中文](https://github.com/wzshiming/gen/blob/master/README_cn.md)

## Supported

- [X] Generate documentation
  - [X] [OpenAPI 3](https://github.com/OAI/OpenAPI-Style-Guide)
  - [X] [SwaggerUI](https://github.com/swagger-api/swagger-ui)
  - [X] [ReDoc](https://github.com/Rebilly/ReDoc)
- [X] RESTful
  - [X] Generate Go router
    - [X] Security
      - [X] apiKey
      - [X] http
        - [X] basic
        - [ ] bearer
      - [ ] oauth2
      - [ ] openIdConnet
    - [X] Content
      - [X] Query
      - [X] Path
      - [X] Header
      - [X] Cookie
      - [X] Body
        - [X] JSON
        - [ ] XML
        - [ ] Formdata
          - [X] File
          - [ ] Value
        - [ ] URLEncode
  - [X] Generate Go client
    - [X] Security
      - [X] apiKey
      - [X] http
        - [X] basic
        - [X] bearer
      - [ ] oauth2
      - [ ] openIdConnet
    - [X] Content
      - [X] Query
      - [X] Path
      - [X] Header
      - [X] Cookie
      - [X] Body
        - [X] JSON
        - [X] XML
        - [X] Formdata
          - [X] File
          - [X] Value
        - [ ] URLEncode
  - [Javascript client](https://github.com/swagger-api/swagger-js)
  - [Other language client](https://github.com/swagger-api/swagger-codegen/tree/3.0.0)
- [ ] gRPC & Proto3

## Examples

1. Install gen tool `go get -v github.com/wzshiming/gen/cmd/gen`
2. Add gen tool to $PATH
3. Start it `gen run github.com/wzshiming/gen-examples/service/...`
4. Open [http://127.0.0.1:8080/swagger/?url=./openapi.json#](http://127.0.0.1:8080/swagger/?url=./openapi.json#) with your browser

[Examples](https://github.com/wzshiming/gen-examples/)  

Or try to quickly build services from scratch

1. Make a directory `mkdir -p $(go env GOPATH)/src/gentest`
2. Change directory `cd $(go env GOPATH)/src/gentest/`
3. Define models
``` shell
cat > models.go <<EOF
package gentest
type Gentest struct {
    Name string \`json:"name"\`
    Age  int    \`json:"age"\`
}
EOF
```
4. Generated from CRUD template `gen crud -t mock -n Gentest`
5. Start it `gen run gentest`

## License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/gen/blob/master/LICENSE) for the full license text.
