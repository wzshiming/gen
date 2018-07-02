package client

import (
	"github.com/wzshiming/gen/client"
	"github.com/wzshiming/gen/parser"
	cli "gopkg.in/urfave/cli.v2"
)

var Command = &cli.Command{
	Name:      "client",
	Usage:     "Generate client source code for functions",
	ArgsUsage: "[package]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Value:   "client_gen.go",
			Usage:   "output file.",
		},
	},
	Action: func(ctx *cli.Context) error {
		pkg := ctx.Args().First()
		out := ctx.String("out")
		if pkg == "" {
			return cli.ShowSubcommandHelp(ctx)
		}

		def := parser.NewParser()
		err := def.Import(pkg)
		if err != nil {
			return err
		}
		d, err := client.NewGenClient(def.API()).Generate()
		if err != nil {
			return err
		}

		return d.WithFilename(out).Save()
	},
}
