package indicators

import (
	"jarvis-trading-bot/structs"
)

type Indicator interface {
	Name() string
	CalcIndicator(candle *structs.Candlestick) *structs.AnalysisReturn
}
