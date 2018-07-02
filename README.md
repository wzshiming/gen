# Gen - Generated source code tool for micro services

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

- [X] Generate router
- [X] Generate client
  - [X] Golang
  - [ ] JavaScript (For Web)
  - [ ] Java (For Android)
  - [ ] Swift (For iOS)
- [X] Generate documentation
  - [ ] Swagger 2
  - [X] [OpenAPI 3](https://github.com/OAI/OpenAPI-Style-Guide)
  - [X] [SwaggerUI](https://github.com/swagger-api/swagger-ui)
- [X] Protocol
  - [X] HTTP
  - [ ] Protobuf

## Examples

[Route Example](https://github.com/wzshiming/gen/blob/master/examples/route1/)  
[Client Example](https://github.com/wzshiming/gen/blob/master/examples/client1/)  

## License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/gen/blob/master/LICENSE) for the full license text.
