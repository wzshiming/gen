package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/wzshiming/gen/client"
	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	"github.com/wzshiming/gen/ui/swaggerui"
	"github.com/wzshiming/openapi/util"
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
				&cli.StringFlag{
					Name:    "format",
					Aliases: []string{"f"},
					Value:   "json",
					Usage:   "It has to be json or yaml",
				},
				&cli.BoolFlag{
					Name:    "ui",
					Aliases: []string{"u"},
					Usage:   "Show the API web page",
				},
			},
			Action: func(ctx *cli.Context) error {
				p := ctx.String("package")
				s := ctx.StringSlice("servers")
				o := ctx.String("out")
				f := ctx.String("format")

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
				switch f {
				case "json":

				case "yaml":
					d, err = util.JSON2YAML(d)
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("undefined format %s", f)
				}
				if o != "" {
					err := ioutil.WriteFile(o, d, 0666)
					if err != nil {
						return err
					}
				}

				if ctx.Bool("ui") {

					mux := &http.ServeMux{}

					mux.Handle("/", swaggerui.Handle)

					mux.HandleFunc("/openapi."+f, func(w http.ResponseWriter, r *http.Request) {
						http.ServeContent(w, r, "openapi."+f, time.Time{}, bytes.NewReader(d))
					})

					fmt.Printf("Open http://127.0.0.1:8080/?url=openapi.%s# with your browser.\n", f)
					return http.ListenAndServe(":8080", mux)
				}

				fmt.Println(string(d))

				return nil
			},
		},
	},
}
