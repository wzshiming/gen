package run

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/run"
)

var (
	port uint
	way  string
)

func init() {
	flag := Cmd.Flags()
	flag.UintVarP(&port, "port", "p", 8080, "listening port")
	flag.StringVarP(&way, "way", "w", "", "way to export")
}

var Cmd = &cobra.Command{
	Use:   "run [flags] package [package ...]",
	Short: "Run package",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Miss package path")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.Run(args, fmt.Sprintf(":%d", port), way)
	},
}
