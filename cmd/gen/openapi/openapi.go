package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/ui/swaggerui"
	"github.com/wzshiming/openapi/util"
	cli "gopkg.in/urfave/cli.v2"
)

var Command = &cli.Command{
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
			Value:   "./openapi.json",
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
		if p == "" {
			return cli.ShowAppHelp(ctx)
		}

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
		} else {
			fmt.Println(string(d))
		}

		if ctx.Bool("ui") {

			mux := &http.ServeMux{}

			mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handle))

			mux.HandleFunc("/swagger/openapi."+f, func(w http.ResponseWriter, r *http.Request) {
				http.ServeContent(w, r, "openapi."+f, time.Time{}, bytes.NewReader(d))
			})
			fmt.Printf("Open http://127.0.0.1:8080/swagger/?url=./openapi.%s# with your browser.\n", f)
			return http.ListenAndServe(":8080", mux)
		}

		return nil
	},
}
