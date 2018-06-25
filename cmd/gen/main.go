package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wzshiming/gen/client"
	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
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
		{
			Name: "client",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "package",
					Aliases: []string{"p"},
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
				},
			},
			Action: func(ctx *cli.Context) error {
				p := ctx.String("package")
				o := ctx.String("out")

				def := parser.NewParser()
				err := def.Import(p)
				if err != nil {
					return err
				}
				d, err := client.NewGenClient(def.API()).Generate()
				if err != nil {
					return err
				}
				if o == "" {
					fmt.Println(string(d))
					return nil
				}
				return ioutil.WriteFile(o, d, 0666)
			},
		},
		{
			Name: "route",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "package",
					Aliases: []string{"p"},
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
				},
			},
			Action: func(ctx *cli.Context) error {
				p := ctx.String("package")
				o := ctx.String("out")

				def := parser.NewParser()
				err := def.Import(p)
				if err != nil {
					return err
				}
				d, err := route.NewGenRoute(def.API()).Generate()
				if err != nil {
					return err
				}
				if o == "" {
					fmt.Println(string(d))
					return nil
				}
				return ioutil.WriteFile(o, d, 0666)
			},
		},
		{
			Name: "openapi",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "package",
					Aliases: []string{"p"},
				},
				&cli.StringSliceFlag{
					Name:    "servers",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
				},
			},
			Action: func(ctx *cli.Context) error {
				p := ctx.String("package")
				s := ctx.StringSlice("servers")
				o := ctx.String("out")

				def := parser.NewParser()
				err := def.Import(p)
				if err != nil {
					return err
				}
				api, err := openapi.NewGenOpenAPI(def.API()).WithServices(s...).Generate()
				if err != nil {
					return err
				}
				d, err := json.MarshalIndent(api, "", " ")
				if err != nil {
					return err
				}
				if o == "" {
					fmt.Println(string(d))
					return nil
				}
				return ioutil.WriteFile(o, d, 0666)
			},
		},
	},
}
