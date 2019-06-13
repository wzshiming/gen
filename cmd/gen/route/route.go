package route

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	"github.com/wzshiming/gen/utils"
	oaspec "github.com/wzshiming/openapi/spec"
	"github.com/wzshiming/openapi/util"
)

var (
	out      string
	name     string
	pack     string
	servers  []string
	info     string
	openapiF bool
	way      string
	explode  bool
)

func init() {
	flag := Cmd.Flags()
	flag.StringVarP(&out, "out", "o", "router_gen.go", "output file name")
	flag.StringVarP(&name, "name", "n", "Router", "routing function name")
	flag.StringVarP(&pack, "package", "p", "", "package name")
	flag.BoolVarP(&openapiF, "openapi", "", false, "with openapi")
	flag.StringSliceVarP(&servers, "servers", "s", nil, "")
	flag.StringVarP(&info, "info", "i", "", "Info")
	flag.StringVarP(&way, "way", "w", "", "way to export")
	flag.BoolVarP(&explode, "explode", "", false, "query parameter of array type explode")
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
		dir, _ := filepath.Abs(out)
		dir = filepath.Dir(dir)
		impPath := utils.GetPackagePath(dir)
		def := parser.NewParser(nil)
		for _, arg := range args {
			err := def.Import(arg, way)
			if err != nil {
				return err
			}
		}

		rg := route.NewGenRoute(def.API()).
			SetExplode(explode)
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

			api, err := openapi.NewGenOpenAPI(def.API()).
				WithServices(servers...).
				SetInfo(oainfo).
				SetExplode(explode).
				Generate()
			if err != nil {
				return err
			}
			rg = rg.WithOpenAPI(api)
		}
		d, err := rg.Generate(pack, impPath, name)
		if err != nil {
			return err
		}

		return d.WithFilename(out).Save()
	},
}
