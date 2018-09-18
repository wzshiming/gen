# Gen - 为微服务生成源码的工具

只需要写普通的函数,Gen会为其生成高效的路由源代码和文档  
因为生成的是源码，所以这些都不会影响运行时性能  
工具中的每个变化所引起的差异直接显示在生成的源代码中  
也支持生成客户端  

[![Build Status](https://travis-ci.org/wzshiming/gen.svg?branch=master)](https://travis-ci.org/wzshiming/gen)
[![Go Report Card](https://goreportcard.com/badge/github.com/wzshiming/gen)](https://goreportcard.com/report/github.com/wzshiming/gen)
[![GitHub license](https://img.shields.io/github/license/wzshiming/gen.svg)](https://github.com/wzshiming/gen/blob/master/LICENSE)

- [English](https://github.com/wzshiming/gen/blob/master/README.md)
- [简体中文](https://github.com/wzshiming/gen/blob/master/README_cn.md)

## 支持的功能

- [X] 生成路由
- [X] 生成客户端
  - [X] Golang
  - [ ] JavaScript
- [X] 生成文档
  - [ ] Swagger 2
  - [X] [OpenAPI 3](https://github.com/OAI/OpenAPI-Style-Guide)
  - [X] [SwaggerUI](https://github.com/swagger-api/swagger-ui)
- [X] 协议
  - [X] HTTP
  - [ ] Protobuf

## 示例

1. 安装 gen 工具 `go get -v github.com/wzshiming/gen/cmd/gen`
2. 添加 gen 工具到 $PATH
3. 执行 `gen run github.com/wzshiming/gen/examples/basics/service`
4. 在浏览器中打开 [http://127.0.0.1:8080/swagger/?url=./openapi.json#](http://127.0.0.1:8080/swagger/?url=./openapi.json#)

[示例](https://github.com/wzshiming/gen/blob/master/examples/)  

## 许可证

软包根据MIT License。有关完整的许可证文本，请参阅[LICENSE](https://github.com/wzshiming/gen/blob/master/LICENSE)。