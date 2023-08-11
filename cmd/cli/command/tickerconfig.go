package command

import (
	"jarvis-trading-bot/cmd/cli/helpers"

	"github.com/spf13/cobra"
)

var (
	cmdTickerConfig = &cobra.Command{
		Use:   "ticker",
		Short: "Manage Tickers configuration in the API",
	}

	cmdCreateTickerConfig = &cobra.Command{
		Use:   "create",
		Short: "Create a new Ticker to use in the API",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helpers.CreateTickerConfig(args[0])
		},
	}

	cmdDeleteTickerConfig = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Ticker Config in the API",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			helpers.DeleteTickerConfig(args[0])
		},
	}

	cmdListTickerConfig = &cobra.Command{
		Use:   "list",
		Short: "List the available Tickers in the API",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			helpers.ListTickers()
		},
	}
)

func init() {
	cmdTickerConfig.AddCommand(cmdCreateTickerConfig)
	cmdTickerConfig.AddCommand(cmdDeleteTickerConfig)
	cmdTickerConfig.AddCommand(cmdListTickerConfig)
	rootCmd.AddCommand(cmdTickerConfig)
}
