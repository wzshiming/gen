package main

import (
	"github.com/spf13/cobra"
	"github.com/wzshiming/gen/cmd/gen/client"
	"github.com/wzshiming/gen/cmd/gen/openapi"
	"github.com/wzshiming/gen/cmd/gen/route"
	"github.com/wzshiming/gen/cmd/gen/run"
)

func main() {
	rootCmd.AddCommand(client.Cmd, run.Cmd, route.Cmd, openapi.Cmd)
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "generated source code tool for micro services",
}
