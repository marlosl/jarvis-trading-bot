package pricerequest

import (
	"fmt"

	"jarvis-trading-bot/services/broker"
	"jarvis-trading-bot/services/exchangeconfig"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/utils"
)

func DoPriceRequest(signalItem *signal.SignalItem) error {
	api := broker.NewBinanceApi()

	symbol := signalItem.Ticker
	exchange := signalItem.Exchange

	exConfig, err := exchangeconfig.GetExchangeConfig(symbol, *exchange)
	if err != nil {
		fmt.Printf("Can't get exchange config: %v", err)
		return err
	}

	if len(exConfig.Symbol) == 0 {
		symbol = exConfig.Symbol
	}

	candle, err := api.GetCandlestick(symbol)
	if err != nil {
		fmt.Printf("Can't get candlestick: %v", err)
		return err
	}

	signalItem.High = &candle.HighPrice
	signalItem.Low = &candle.LowPrice
	signalItem.Open = &candle.OpenPrice
	signalItem.Close = &candle.ClosePrice
	signalItem.Volume = &candle.Volume
	signalItem.Time = &candle.KlineCloseTime
	signalItem.TimeNow = &candle.KlineCloseTime
	signalItem.TakeProfitPerc = &exConfig.TakeProfitPerc
	signalItem.StopLossPerc = &exConfig.StopLossPerc

	fmt.Println("Price response: ", utils.SPrintJson(signalItem))
	return nil
}
