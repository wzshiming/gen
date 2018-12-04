package openapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/ui/swaggerui"
	"github.com/wzshiming/gotype"
	"github.com/wzshiming/openapi/util"
)

var (
	servers []string
	port    uint
	out     string
	format  string
	ui      bool
)

func init() {
	flag := Cmd.Flags()
	flag.StringSliceVarP(&servers, "servers", "s", nil, "")
	flag.UintVarP(&port, "port", "p", 8080, "listening port")
	flag.StringVarP(&out, "out", "o", "openapi.json", "output file name")
	flag.StringVarP(&format, "format", "f", "json", "json or yaml")
	flag.BoolVarP(&ui, "ui", "u", false, "show the API ui page")

}

var Cmd = &cobra.Command{
	Use:   "openapi [flags] package",
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
		err := def.Import(args[0])
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

		if ui {

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
