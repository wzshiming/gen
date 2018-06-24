package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/wzshiming/ffmt"
	"github.com/wzshiming/gen/api"
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

	def := api.NewGenAPI()
	err := def.Import("github.com/wzshiming/gen/testdata")
	if err != nil {
		ffmt.Mark(err)
		return
	}

	swa := def.OpenAPI()
	writeJson(swa)
}
