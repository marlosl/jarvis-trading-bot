package analyser

import (
	"testing"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/services/types"
	"jarvis-trading-bot/utils"

	"github.com/stretchr/testify/assert"
)

var transactionItem = &types.TransactionItem{
	Uuid:      "22412091-cdaf-47a8-8b89-144d8009b9e4",
	CreatedAt: nil,
	Exchange:  utils.NewString("BINANCE"),
	Ticker:    "BTCUSD",
	Status:    consts.STATUS_ACTIVE,
	Signals: []types.Signal{
		{
			SignalTime:    utils.GetCurrentTime(),
			IndicatorName: signal.ALGOD,
			Action:        signal.BUY,
		},
		{
			SignalTime:    utils.GetCurrentTime(),
			IndicatorName: signal.ALGOD,
			Action:        signal.BUY,
		},
		{
			SignalTime:    utils.GetCurrentTime(),
			IndicatorName: signal.ALGOD,
			Action:        "TRENDCHANGED",
		},
		{
			SignalTime:    utils.GetCurrentTime(),
			IndicatorName: signal.ALGOD,
			Action:        "TRENDCHANGED",
		},
	},
}

func TestIfChangeToClosedStatus(t *testing.T) {
	assert := assert.New(t)
	analyser := NewAnalyserService()

	signalItem := &signal.SignalItem{
		Alert:         "ALGOD_SELL",
		Exchange:      utils.NewString("BINANCE"),
		Ticker:        "BTCUSD",
		Action:        signal.SELL,
		IndicatorName: signal.ALGOD,
		Close:         utils.NewString("24475.82"),
		Open:          utils.NewString("24494.54"),
		High:          utils.NewString("24507.13"),
		Low:           utils.NewString("24475.82"),
		Time:          utils.NewString("2023-02-17T22:36:00Z"),
		Volume:        utils.NewString("4.71386"),
		TimeNow:       utils.NewString("2023-02-17T22:36:58Z"),
	}

	item := transactionItem

	item = analyser.DoAnalysis(item, signalItem)
	assert.Equal(item.Status, consts.STATUS_CLOSED)
}

func TestIfTransactionItemIsNotNull(t *testing.T) {
	assert := assert.New(t)
	analyser := NewAnalyserService()

	signalItem := &signal.SignalItem{
		Alert:         "ALGOD_BUY",
		Exchange:      utils.NewString("BINANCE"),
		Ticker:        "BTCUSD",
		Action:        signal.BUY,
		IndicatorName: signal.ALGOD,
		Close:         utils.NewString("24475.82"),
		Open:          utils.NewString("24494.54"),
		High:          utils.NewString("24507.13"),
		Low:           utils.NewString("24475.82"),
		Time:          utils.NewString("2023-02-17T22:36:00Z"),
		Volume:        utils.NewString("4.71386"),
		TimeNow:       utils.NewString("2023-02-17T22:36:58Z"),
	}

	item := analyser.DoAnalysis(nil, signalItem)
	assert.NotNil(item)
}

func TestIfTransactionItemIsNull(t *testing.T) {
	assert := assert.New(t)
	analyser := NewAnalyserService()

	signalItem := &signal.SignalItem{
		Alert:         "ALGOD_SELL",
		Exchange:      utils.NewString("BINANCE"),
		Ticker:        "BTCUSD",
		Action:        signal.SELL,
		IndicatorName: signal.ALGOD,
		Close:         utils.NewString("24475.82"),
		Open:          utils.NewString("24494.54"),
		High:          utils.NewString("24507.13"),
		Low:           utils.NewString("24475.82"),
		Time:          utils.NewString("2023-02-17T22:36:00Z"),
		Volume:        utils.NewString("4.71386"),
		TimeNow:       utils.NewString("2023-02-17T22:36:58Z"),
	}

	item := analyser.DoAnalysis(nil, signalItem)
	assert.Nil(item)
}

func TestStopLoss(t *testing.T) {
	assert := assert.New(t)
	analyser := NewAnalyserService()

	signalItem := &signal.SignalItem{
		Alert:         "",
		Exchange:      utils.NewString("BINANCE"),
		Ticker:        "BTCUSD",
		Action:        signal.PRICE_REQUEST,
		IndicatorName: "",
		Close:         utils.NewString("24475.82"),
		Open:          utils.NewString("24494.54"),
		High:          utils.NewString("24507.13"),
		Low:           utils.NewString("24475.82"),
		Time:          utils.NewString("2023-02-17T22:36:00Z"),
		Volume:        utils.NewString("4.71386"),
		StopLossPerc:  utils.NewString("10"),
		TimeNow:       utils.NewString("2023-02-17T22:36:58Z"),
	}

	item := transactionItem
	item.BuyPrice = utils.NewString("28000.82")

	item = analyser.DoAnalysis(item, signalItem)

	hasStopLoss := false
	for _, sig := range item.Signals {
		if sig.Action == signal.STOP_LOSS {
			hasStopLoss = true
		}
	}

	assert.True(hasStopLoss)
	assert.NotNil(item.SellPrice)
	assert.Equal(item.Status, consts.STATUS_CLOSED)
}

func TestTakeProfit(t *testing.T) {
	assert := assert.New(t)
	analyser := NewAnalyserService()

	signalItem := &signal.SignalItem{
		Alert:          "",
		Exchange:       utils.NewString("BINANCE"),
		Ticker:         "BTCUSD",
		Action:         signal.PRICE_REQUEST,
		IndicatorName:  "",
		Close:          utils.NewString("28000.82"),
		Open:           utils.NewString("24494.54"),
		High:           utils.NewString("24507.13"),
		Low:            utils.NewString("24475.82"),
		Time:           utils.NewString("2023-02-17T22:36:00Z"),
		Volume:         utils.NewString("4.71386"),
		TakeProfitPerc: utils.NewString("10"),
		TimeNow:        utils.NewString("2023-02-17T22:36:58Z"),
	}

	item := transactionItem
	item.BuyPrice = utils.NewString("24475.82")

	item = analyser.DoAnalysis(item, signalItem)

	hasTakeProfit := false
	for _, sig := range item.Signals {
		if sig.Action == signal.TAKE_PROFIT {
			hasTakeProfit = true
		}
	}

	assert.True(hasTakeProfit)
	assert.NotNil(item.SellPrice)
	assert.Equal(item.Status, consts.STATUS_CLOSED)
}
