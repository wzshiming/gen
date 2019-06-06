package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	oaspec "github.com/wzshiming/openapi/spec"
	"github.com/wzshiming/openapi/util"
)

var (
	servers []string
	out     string
	format  string
	info    string
	way     string
)

func init() {
	flag := Cmd.Flags()
	flag.StringSliceVarP(&servers, "servers", "s", nil, "")
	flag.StringVarP(&out, "out", "o", "openapi.json", "output file name")
	flag.StringVarP(&format, "format", "f", "json", "json or yaml")
	flag.StringVarP(&info, "info", "i", "", "Info")
	flag.StringVarP(&way, "way", "w", "", "way to export")
}

var Cmd = &cobra.Command{
	Use:   "openapi [flags] package [package ...]",
	Short: "Generate openapi document for functions",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Miss package path")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		def := parser.NewParser(nil)
		for _, arg := range args {
			err := def.Import(arg, way)
			if err != nil {
				return err
			}
		}

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

		return nil

	},
}
