package run

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/run"
)

var (
	port   uint
	format string
)

func init() {
	flag := Cmd.Flags()
	flag.UintVarP(&port, "port", "p", 8080, "listening port")
	flag.StringVarP(&format, "format", "f", "json", "json or yaml")
}

var Cmd = &cobra.Command{
	Use:   "run [flags] package",
	Short: "Run package",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Miss package path")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.Run(args[0], fmt.Sprintf(":%d", port), format)
	},
}
