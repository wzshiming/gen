package run

import (
	"fmt"

	"github.com/wzshiming/gen/run"
	cli "gopkg.in/urfave/cli.v2"
)

var Command = &cli.Command{
	Name:      "run",
	Usage:     "Run package",
	ArgsUsage: "[package]",
	Flags: []cli.Flag{
		&cli.UintFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   8080,
			Usage:   "listening addrs.",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Value:   "json",
			Usage:   "json or yaml.",
		},
	},
	Action: func(ctx *cli.Context) error {
		pkg := ctx.Args().First()
		port := ctx.Uint("port")
		format := ctx.String("format")
		if pkg == "" {
			return cli.ShowSubcommandHelp(ctx)
		}
		return run.Run(pkg, fmt.Sprintf(":%d", port), format)
	},
}
