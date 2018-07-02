package main

import (
	"os"

	"github.com/wzshiming/gen/cmd/gen/client"
	"github.com/wzshiming/gen/cmd/gen/openapi"
	"github.com/wzshiming/gen/cmd/gen/route"
	"github.com/wzshiming/gen/cmd/gen/run"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(2)
		return
	}
}

var app = &cli.App{
	Name:    "gen",
	Version: "0.0.X",
	Authors: []*cli.Author{
		{
			Name:  "wzshiming",
			Email: "wzshiming@foxmail.com",
		},
	},
	Usage:     "generated source code tool for micro services",
	ArgsUsage: "[package]",
	Commands: []*cli.Command{
		client.Command,
		route.Command,
		openapi.Command,
		run.Command,
	},
}
