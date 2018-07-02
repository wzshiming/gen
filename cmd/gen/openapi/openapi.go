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
	Name:      "openapi",
	Usage:     "Generate openapi document for functions",
	ArgsUsage: "[package]",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "servers",
			Aliases: []string{"s"},
		},
		&cli.UintFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   8080,
			Usage:   "listening addrs.",
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Value:   "./openapi.json",
			Usage:   "output file.",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Value:   "json",
			Usage:   "json or yaml.",
		},
		&cli.BoolFlag{
			Name:    "ui",
			Aliases: []string{"u"},
			Usage:   "show the API web page.",
		},
	},
	Action: func(ctx *cli.Context) error {
		pkg := ctx.Args().First()
		servers := ctx.StringSlice("servers")
		out := ctx.String("out")
		format := ctx.String("format")
		port := ctx.Uint("port")
		if pkg == "" {
			return cli.ShowSubcommandHelp(ctx)
		}

		def := parser.NewParser()
		err := def.Import(pkg)
		if err != nil {
			return err
		}
		api, err := openapi.NewGenOpenAPI(def.API()).WithServices(servers...).Generate()
		if err != nil {
			return err
		}

		d, err := json.MarshalIndent(api, "", " ")
		if err != nil {
			return err
		}
		switch format {
		case "json":

		case "yaml":
			d, err = util.JSON2YAML(d)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("undefined format %s", format)
		}
		if out != "" {
			err := ioutil.WriteFile(out, d, 0666)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(string(d))
		}

		if ctx.Bool("ui") {

			mux := &http.ServeMux{}

			mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handle))

			mux.HandleFunc("/swagger/openapi."+format, func(w http.ResponseWriter, r *http.Request) {
				http.ServeContent(w, r, "openapi."+format, time.Time{}, bytes.NewReader(d))
			})
			fmt.Printf("Open http://127.0.0.1:%d/swagger/?url=./openapi.%s# with your browser.\n", port, format)
			return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
		}

		return nil
	},
}
