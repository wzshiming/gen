package client

import (
	"github.com/wzshiming/gen/client"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gotype"
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
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"p"},
			Usage:   "package name",
		},
	},
	Action: func(ctx *cli.Context) error {
		pkg := ctx.Args().First()
		out := ctx.String("out")
		packag := ctx.String("package")
		if pkg == "" {
			return cli.ShowSubcommandHelp(ctx)
		}

		imp := gotype.NewImporter(gotype.WithCommentLocator())
		def := parser.NewParser(imp)
		err := def.Import(pkg)
		if err != nil {
			return err
		}
		d, err := client.NewGenClient(def.API()).Generate()
		if err != nil {
			return err
		}

		return d.WithPackname(packag).WithFilename(out).Save()
	},
}
