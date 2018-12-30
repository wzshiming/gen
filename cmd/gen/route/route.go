package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/gotype"
	oaspec "github.com/wzshiming/openapi/spec"
	"github.com/wzshiming/openapi/util"
)

var (
	out      string
	name     string
	pack     string
	servers  []string
	format   string
	info     string
	openapiF bool
)

func init() {
	flag := Cmd.Flags()
	flag.StringVarP(&out, "out", "o", "router_gen.go", "output file name")
	flag.StringVarP(&name, "name", "n", "Router", "routing function name")
	flag.StringVarP(&pack, "package", "p", "", "package name")
	flag.BoolVarP(&openapiF, "openapi", "", false, "with openapi")
	flag.StringSliceVarP(&servers, "servers", "s", nil, "")
	flag.StringVarP(&format, "format", "f", "json", "json or yaml")
	flag.StringVarP(&info, "info", "i", "", "Info")
}

var Cmd = &cobra.Command{
	Use:   "route [flags] package [package ...]",
	Short: "Generate routing source code for functions",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Miss package path")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		imp := gotype.NewImporter(gotype.WithCommentLocator())
		dir, _ := filepath.Abs(out)
		dir = filepath.Dir(dir)
		impPath := utils.GetPackagePath(dir)
		def := parser.NewParser(imp)
		for _, arg := range args {
			err := def.Import(arg)
			if err != nil {
				return err
			}
		}
		d, err := route.NewGenRoute(def.API()).Generate(pack, impPath, name)
		if err != nil {
			return err
		}

		if openapiF {

			var oainfo *oaspec.Info

			if info != "" {
				fil, err := ioutil.ReadFile(info)
				if err != nil {
					return err
				}
				fil, err = util.YAML2JSON(fil)

				if err != nil {
					return err
				}
				err = json.Unmarshal(fil, &oainfo)
				if err != nil {
					return err
				}
			}

			api, err := openapi.NewGenOpenAPI(def.API()).WithServices(servers...).SetInfo(oainfo).Generate()
			if err != nil {
				return err
			}

			dc, err := json.MarshalIndent(api, "", " ")
			if err != nil {
				return err
			}
			switch format {
			case "json":

			case "yaml":
				dc, err = util.JSON2YAML(dc)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("undefined format %s", format)
			}

			d.WriteFormat("var OpenAPI=`%s`", string(dc))
		}

		return d.WithFilename(out).Save()
	},
}
