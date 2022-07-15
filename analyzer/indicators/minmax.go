package indicators

import (
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/utils/log"

	"github.com/shopspring/decimal"
)

const (
	MinMaxIndicator = "MinMax"
)

type MinMax struct {
	Crypto  structs.BotParams
	Candles []structs.Candlestick
}

func (m *MinMax) Name() string {
	return MinMaxIndicator
}

func (m *MinMax) CalcIndicator(candle *structs.Candlestick) *structs.AnalysisReturn {
	log.InfoLogger.Printf("Symbol: %s\n", m.Crypto.Symbol)
	m.doCalcMinMax(m.Crypto.Symbol)
	log.InfoLogger.Println("")
	return nil
}

func (m *MinMax) doCalcMinMax(symbol string) {
	startValue := decimal.NewFromFloat(0)
	minValue := decimal.NewFromFloat(0)
	maxValue := decimal.NewFromFloat(0)
	for i, candle := range m.Candles {
		if i == 0 {
			startValue = candle.ClosePrice
			minValue = candle.ClosePrice
			maxValue = candle.ClosePrice
			continue
		}

		if candle.ClosePrice.Cmp(maxValue) > 0 {
			maxValue = candle.ClosePrice
		}

		if candle.ClosePrice.Cmp(minValue) < 0 {
			minValue = candle.ClosePrice
		}
	}

	minDiff := startValue.Sub(minValue)
	minVar := minDiff.Div(startValue).Mul(decimal.NewFromInt(100))

	maxDiff := maxValue.Sub(startValue)
	maxVar := maxDiff.Div(startValue).Mul(decimal.NewFromInt(100))

	log.InfoLogger.Printf("startValue: %s\n", startValue)
	log.InfoLogger.Printf("minValue: %s, minVar: %s\n", minValue, minVar)
	log.InfoLogger.Printf("maxValue: %s, maxVar: %s\n", maxValue, maxVar)
}
