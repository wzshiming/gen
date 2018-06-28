package client

import (
	"fmt"
	"io/ioutil"

	"github.com/wzshiming/gen/client"
	"github.com/wzshiming/gen/parser"
	cli "gopkg.in/urfave/cli.v2"
)

var Command = &cli.Command{
	Name: "client",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Value:   "client_gen.go",
		},
	},
	Action: func(ctx *cli.Context) error {
		p := ctx.String("package")
		o := ctx.String("out")
		if p == "" {
			return cli.ShowAppHelp(ctx)
		}

		def := parser.NewParser()
		err := def.Import(p)
		if err != nil {
			return err
		}
		d, err := client.NewGenClient(def.API()).Generate()
		if err != nil {
			return err
		}
		if o == "" {
			fmt.Println(string(d))
			return nil
		}
		return ioutil.WriteFile(o, d, 0666)
	},
}
