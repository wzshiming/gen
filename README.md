# Gen - RESTful style of micro services to generate source code tools

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
- [X] RESTful
  - [X] Generate Go router
  - [X] Generate Go client
  - [Javascript client](https://github.com/swagger-api/swagger-js)
  - [Other language client](https://github.com/swagger-api/swagger-codegen)
- [ ] gRPC & Proto3

## Examples

1. Install gen tool `go get -v github.com/wzshiming/gen/cmd/gen`
2. Add gen tool to $PATH
3. Execute it `gen run github.com/wzshiming/gen/examples/basics/service`
4. Open [http://127.0.0.1:8080/swagger/?url=./openapi.json#](http://127.0.0.1:8080/swagger/?url=./openapi.json#) with your browser

[Examples](https://github.com/wzshiming/gen/blob/master/examples/)  

## License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/gen/blob/master/LICENSE) for the full license text.
