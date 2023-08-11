package command

import (
	"jarvis-trading-bot/cmd/cli/helpers"
	"strings"

	"github.com/spf13/cobra"
)

var (
	buildBinariesCmd = &cobra.Command{
		Use:       "build-binaries",
		Short:     "Build the binaries for the Jarvis Trading Bot.",
		Long:      "Build the binaries for the Jarvis Trading Bot.\nValid options are: " + strings.Join(helpers.GetFunctionNames(), ", ") + ". Use no options will build all binaries.",
		ValidArgs: helpers.GetFunctionNames(),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				helpers.BuildBinaryFiles()
				return
			}

			helpers.BuildSingleBinaryFile(args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(buildBinariesCmd)
}
