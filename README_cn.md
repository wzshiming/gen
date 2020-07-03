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

## 示例

'#'是注释，注释是golang tag语法，这里唯一的区别是'#' 包裹而不是 '`'

``` golang
// ItemService #path:"/item/"#
type ItemService struct {}

// Create a Item #route:"POST /"#
func (s *ItemService) Create(item *Item) (err error) {}

// Update the Item #route:"PUT /{item_id}"#
func (s *ItemService) Update(itemID int /* #name:"item_id"# */, item *Item) (err error) {}

// Delete the Item #route:"DELETE /{item_id}"#
func (s *ItemService) Delete(itemID int /* #name:"item_id"# */) (err error) {}

// Get the Item #route:"GET /{item_id}"#
func (s *ItemService) Get(itemID int /* #name:"item_id"# */) (item *ItemWithID, err error) {}

// List of the Item #route:"GET /"#
func (s *ItemService) List(offset, limit int) (items []*ItemWithID, err error) {}
```

1. 安装 gen 工具 `go get -v github.com/wzshiming/gen/cmd/gen`
2. 添加 gen 工具到 $PATH
3. 启动 `gen run github.com/wzshiming/gen-examples/service/...`
4. 在浏览器中打开 [http://127.0.0.1:8080/swagger/?url=./openapi.json#](http://127.0.0.1:8080/swagger/?url=./openapi.json#)

[示例](https://github.com/wzshiming/gen-examples/)  

或者尝试从零快速搭建web服务

1. 新建目录 `mkdir -p $(go env GOPATH)/src/gentest`
2. 移动到刚创建的目录 `cd $(go env GOPATH)/src/gentest/`
3. 定义数据类型
``` shell
cat > models.go <<EOF
package gentest
type Gentest struct {
    Name string \`json:"name"\`
    Age  int    \`json:"age"\`
}
EOF
```
4. 根据CRUD模板生成
`gen crud -t mock -n Gentest`
5. 启动
`GO111MODULE=off gen run gentest`

## 支持的功能

- [X] 生成文档
  - [X] [OpenAPI 3](https://github.com/OAI/OpenAPI-Style-Guide)
  - [X] [SwaggerUI](https://github.com/swagger-api/swagger-ui)
  - [X] [ReDoc](https://github.com/Rebilly/ReDoc)
- [X] RESTful
  - [X] 生成Go路由
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
  - [X] 生成Go客户端
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
  - 其他客户端
    - [OpenAPITools/openapi-generator](https://github.com/OpenAPITools/openapi-generator)
    - [swagger-api/swagger-js](https://github.com/swagger-api/swagger-js)
    - [swagger-api/swagger-codegen](https://github.com/swagger-api/swagger-codegen/tree/3.0.0)
- [ ] gRPC & Proto3

## 许可证

软包根据MIT License。有关完整的许可证文本，请参阅[LICENSE](https://github.com/wzshiming/gen/blob/master/LICENSE)。
