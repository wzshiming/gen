package route

import (
	"path"

	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	"github.com/wzshiming/gotype"
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
		name := ctx.String("name")
		if pkg == "" {
			return cli.ShowSubcommandHelp(ctx)
		}
		imp := gotype.NewImporter(gotype.WithCommentLocator())
		impPath := ""
		if outpath, _ := imp.ImportBuild(path.Dir(path.Join(pkg, out))); outpath != nil {
			impPath = outpath.ImportPath
		}
		def := parser.NewParser(imp)
		err := def.Import(pkg)
		if err != nil {
			return err
		}
		d, err := route.NewGenRoute(def.API()).Generate(packag, impPath, name)
		if err != nil {
			return err
		}
		return d.WithFilename(out).Save()
	},
}
