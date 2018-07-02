package route

import (
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	cli "gopkg.in/urfave/cli.v2"
)

var Command = &cli.Command{
	Name:      "route",
	Usage:     "Generate routing source code for functions",
	ArgsUsage: "[package]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Value:   "router_gen.go",
			Usage:   "file name.",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "Router",
			Usage:   "routing function name.",
		},
	},
	Action: func(ctx *cli.Context) error {
		pkg := ctx.Args().First()
		out := ctx.String("out")
		name := ctx.String("name")
		if pkg == "" {
			return cli.ShowSubcommandHelp(ctx)
		}

		def := parser.NewParser()
		err := def.Import(pkg)
		if err != nil {
			return err
		}
		d, err := route.NewGenRoute(def.API()).Generate(name)
		if err != nil {
			return err
		}
		return d.WithFilename(out).Save()
	},
}
