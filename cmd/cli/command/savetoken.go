package command

import (
	"jarvis-trading-bot/services/authentication"

	"github.com/spf13/cobra"
)

var cmdSaveToken = &cobra.Command{
	Use:   "save-token",
	Short: "Save a token to use in the API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authentication.SaveSecret(args[0])
	},
}

func init() {
	rootCmd.AddCommand(cmdSaveToken)
}
