package client

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/client"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gotype"
)

var (
	out  string
	pack string
	way  string
)

func init() {
	flag := Cmd.Flags()
	flag.StringVarP(&out, "out", "o", "client_gen.go", "output file")
	flag.StringVarP(&pack, "package", "p", "", "package name")
	flag.StringVarP(&way, "way", "w", "", "way to export")
}

var Cmd = &cobra.Command{
	Use:   "client [flags] package [package...]",
	Short: "Generate client source code for functions",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Miss package path")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		imp := gotype.NewImporter(gotype.WithCommentLocator())
		def := parser.NewParser(imp)

		for _, pkg := range args {
			err := def.Import(pkg, way)
			if err != nil {
				return err
			}
		}
		d, err := client.NewGenClient(def.API()).Generate()
		if err != nil {
			return err
		}

		return d.WithPackname(pack).WithFilename(out).Save()
	},
}
