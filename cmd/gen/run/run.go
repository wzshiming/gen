package run

import (
	"github.com/wzshiming/gen/run"
	cli "gopkg.in/urfave/cli.v2"
)

var Command = &cli.Command{
	Name: "run",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:    "server",
			Aliases: []string{"s"},
			Value:   ":8080",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Value:   "json",
			Usage:   "It has to be json or yaml for openapi",
		},
	},
	Action: func(ctx *cli.Context) error {
		p := ctx.String("package")
		s := ctx.String("server")
		f := ctx.String("format")
		return run.Run(p, s, f)
	},
}
