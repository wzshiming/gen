package route

import (
	"errors"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	"github.com/wzshiming/gen/utils"
	"github.com/wzshiming/gotype"
)

var (
	out  string
	name string
	pack string
)

func init() {
	flag := Cmd.Flags()
	flag.StringVarP(&out, "out", "o", "router_gen.go", "output file name")
	flag.StringVarP(&name, "name", "n", "Router", "routing function name")
	flag.StringVarP(&pack, "package", "p", "", "package name")

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
		return d.WithFilename(out).Save()
	},
}
