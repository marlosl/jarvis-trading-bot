package command

import (
	"fmt"
	"jarvis-trading-bot/cmd/cli/helpers"

	"github.com/spf13/cobra"
)

var (
	tickers  []string
	filename string

	cmdGetTransactions = &cobra.Command{
		Use:   "get-transactions",
		Short: "Get transactions results",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			if len(tickers) > 0 && len(filename) > 0 {
				helpers.GenerateFile(tickers, filename)
				return
			}

			if len(tickers) > 0 {
				helpers.GetTransactions(tickers)
				return
			}

			fmt.Println("Please, provide a ticker or a filename to save the results\n Use \"ticket list\" to get the available tickers")
		},
	}
)

func init() {
	cmdGetTransactions.Flags().StringSliceVarP(&tickers, "tickers", "t", []string{}, "Tickers to get transactions")
	cmdGetTransactions.Flags().StringVarP(&filename, "filename", "f", "", "CSV filename to save the results")

	rootCmd.AddCommand(cmdGetTransactions)
}
