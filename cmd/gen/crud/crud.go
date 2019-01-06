package crud

import (
	"errors"
	"go/build"
	"go/format"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/crud"
	"github.com/wzshiming/namecase"
)

var (
	pkgname string
	tplname string
	name    string
	out     string
)

func init() {
	flag := Cmd.Flags()
	flag.StringVarP(&tplname, "tpl", "t", "mock", "tpl name (mock, mgo)")
	flag.StringVarP(&name, "name", "n", "", "type name")
	flag.StringVarP(&pkgname, "package", "p", "", "package name")
	flag.StringVarP(&out, "out", "o", "", "out file name")
}

var Cmd = &cobra.Command{
	Use:   "crud [flags]",
	Short: "Generate CRUD source code of type",
	Args: func(cmd *cobra.Command, args []string) error {
		if name == "" {
			return errors.New("need the type name")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if pkgname == "" {
			i, err := build.Import(".", ".", 0)
			if err == nil {
				pkgname = i.Name
			}
		}

		if pkgname == "" {
			pkgname = "main"
		}

		g := crud.NewGenCrud()
		data, err := g.Generate(tplname, pkgname, name)
		if err != nil {
			return err
		}

		if out == "" {
			out = namecase.ToLowerSnake(name) + ".go"
		}

		if data0, err := format.Source(data); err == nil {
			data = data0
		}

		err = ioutil.WriteFile(out, data, 0666)
		if err != nil {
			return err
		}

		return nil

	},
}
