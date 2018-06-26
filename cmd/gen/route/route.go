package route

import (
	"fmt"
	"io/ioutil"

	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	cli "gopkg.in/urfave/cli.v2"
)

var Command = &cli.Command{
	Name: "route",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
		},
	},
	Action: func(ctx *cli.Context) error {
		p := ctx.String("package")
		o := ctx.String("out")

		def := parser.NewParser()
		err := def.Import(p)
		if err != nil {
			return err
		}
		d, err := route.NewGenRoute(def.API()).Generate()
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
