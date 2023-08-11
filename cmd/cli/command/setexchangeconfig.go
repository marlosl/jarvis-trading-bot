package command

import (
	"jarvis-trading-bot/cmd/cli/helpers"

	"github.com/spf13/cobra"
)

var (
	exchange         string
	ticker           string
	symbol           string
	buyQty           string
	stopLossPerc     string
	takeProfitPerc   string
	realTransactions string

	cmdSetExchangeConfig = &cobra.Command{
		Use:   "set-exchange-config",
		Short: "Set exchange configuration",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			helpers.SetExchangeConfig(exchange, ticker, symbol, buyQty, stopLossPerc, takeProfitPerc, realTransactions == "true")
		},
	}
)

func init() {
	cmdSetExchangeConfig.Flags().StringVarP(&exchange, "exchange", "e", "BINANCE", "Exchange to set configuration")
	cmdSetExchangeConfig.Flags().StringVarP(&ticker, "ticker", "t", "BTCUSD", "Ticker to set configuration")
	cmdSetExchangeConfig.Flags().StringVarP(&symbol, "symbol", "s", "BTCUSDT", "Symbol to set configuration")
	cmdSetExchangeConfig.Flags().StringVarP(&buyQty, "buy-qty", "b", "0.01", "Buy quantity to set configuration")
	cmdSetExchangeConfig.Flags().StringVarP(&stopLossPerc, "stop-loss", "l", "0", "Stop Loss Percentage to set configuration")
	cmdSetExchangeConfig.Flags().StringVarP(&takeProfitPerc, "take-profit", "p", "0", "Take Profit Percentage to set configuration")
	cmdSetExchangeConfig.Flags().StringVarP(&realTransactions, "real-transactions", "r", "false", "Real Transactions to set configuration")

	rootCmd.AddCommand(cmdSetExchangeConfig)
}
