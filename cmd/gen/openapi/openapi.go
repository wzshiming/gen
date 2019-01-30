package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gotype"
	oaspec "github.com/wzshiming/openapi/spec"
	oaui "github.com/wzshiming/openapi/ui"
	"github.com/wzshiming/openapi/ui/redoc"
	"github.com/wzshiming/openapi/ui/swaggerui"
	"github.com/wzshiming/openapi/util"
)

var (
	servers []string
	port    uint
	out     string
	format  string
	ui      bool
	info    string
	way     string
)

func init() {
	flag := Cmd.Flags()
	flag.StringSliceVarP(&servers, "servers", "s", nil, "")
	flag.UintVarP(&port, "port", "p", 8080, "listening port")
	flag.StringVarP(&out, "out", "o", "openapi.json", "output file name")
	flag.StringVarP(&format, "format", "f", "json", "json or yaml")
	flag.BoolVarP(&ui, "ui", "u", false, "show the API ui page")
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
		imp := gotype.NewImporter(gotype.WithCommentLocator())
		def := parser.NewParser(imp)
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
		dj := d
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

		if ui {

			mux := &http.ServeMux{}

			mux.Handle("/swagger/", http.StripPrefix("/swagger", oaui.HandleWithFile("openapi.json", dj, swaggerui.Asset)))
			mux.Handle("/redoc/", http.StripPrefix("/redoc", oaui.HandleWithFile("openapi.json", dj, redoc.Asset)))
			fmt.Printf("Open http://127.0.0.1:%d/swagger/# or http://127.0.0.1:%d/redoc/# with your browser.\n", port, port)

			mux0 := handlers.RecoveryHandler()(mux)
			mux0 = handlers.CombinedLoggingHandler(os.Stdout, mux0)
			return http.ListenAndServe(fmt.Sprintf(":%d", port), mux0)
		}
		return nil

	},
}
