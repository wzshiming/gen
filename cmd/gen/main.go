package main

import (
	"os"

	"github.com/wzshiming/gen/cmd/gen/client"
	"github.com/wzshiming/gen/cmd/gen/openapi"
	"github.com/wzshiming/gen/cmd/gen/route"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		return
	}
}

var app = &cli.App{
	Name: "gen",
	Commands: []*cli.Command{
		client.Command,
		route.Command,
		openapi.Command,
	},
}
