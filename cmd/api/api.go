package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/wzshiming/ffmt"
	"github.com/wzshiming/gen"
	"github.com/wzshiming/gen/openapi"
)

func printJson(i interface{}) {
	d, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		ffmt.Mark(err)
		return
	}
	fmt.Println(string(d))
}

func writeJson(i interface{}) {
	d, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		ffmt.Mark(err)
		return
	}
	ffmt.Mark(string(d))
	ioutil.WriteFile("openapi.json", d, 0666)
}

func main() {
	def := gen.NewGen()
	err := def.Import("github.com/wzshiming/gen/testdata")
	if err != nil {
		ffmt.Mark(err)
		return
	}
	api, err := openapi.NewGenOpenAPI(def.API()).WithServices("http://127.0.0.1:8080/").Generate()
	if err != nil {
		ffmt.Mark(err)
		return
	}
	writeJson(api)
}
