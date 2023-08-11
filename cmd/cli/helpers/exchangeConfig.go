package helpers

import (
	"fmt"
	"jarvis-trading-bot/services/exchangeconfig"
)

func SetExchangeConfig(
	exchange string,
	ticker string,
	symbol string,
	buyQty string,
	stopLossPerc string,
	takeProfitPerc string,
	realTransactions bool,
) {
	exchConfig := &exchangeconfig.ExchangeConfigItem{
		Exchange:         exchange,
		Ticker:           ticker,
		Symbol:           symbol,
		BuyQty:           buyQty,
		StopLossPerc:     stopLossPerc,
		TakeProfitPerc:   takeProfitPerc,
		RealTransactions: realTransactions,
	}

	err := exchangeconfig.SaveExchangeConfig(exchConfig)
	if err != nil {
		fmt.Printf("Can't save exchange config: %v", err)
	}
}
