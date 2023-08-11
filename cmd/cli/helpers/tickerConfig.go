package helpers

import (
	"fmt"
	"jarvis-trading-bot/services/tickerconfig"
)

func CreateTickerConfig(ticker string) {
	err := tickerconfig.CreateTickerConfigIfNotExists(ticker)
	if err != nil {
		fmt.Printf("Can't create ticker config: %v", err)
	}
}

func DeleteTickerConfig(ticker string) {
	err := tickerconfig.DeleteTickerConfig(ticker)
	if err != nil {
		fmt.Printf("Delete ticker config: %v", err)
	}
}

func ListTickers() {
	tickers := tickerconfig.ListTickers()

	if len(tickers) == 0 {
		fmt.Println("No Tickers available")
		return
	}

	fmt.Println("")
	fmt.Println("Available Tickers: ")
	fmt.Println("")

	for _, ticker := range tickers {
		fmt.Printf("\t%s\n", ticker)
	}
}
