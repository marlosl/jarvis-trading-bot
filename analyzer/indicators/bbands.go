package indicators

import (
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/talib"

	"github.com/shopspring/decimal"
)

type BBands struct {
	BotParams     structs.BotParams
	TradingStatus *structs.TradingStatus
	Closes        []decimal.Decimal
	OnlyCalculate bool
	Period        int
}

func (b *BBands) CalcBBands(candle *structs.Candlestick) *structs.AnalysisReturn {
	return b.doCalcBBands(candle)
}

func (b *BBands) doCalcBBands(candle *structs.Candlestick) *structs.AnalysisReturn {
	b.Period = 20
	aReturn := new(structs.AnalysisReturn)
	aReturn.Operation = consts.OperationNone
	aReturn.Price = decimal.Zero

	b.Closes = append(b.Closes, candle.ClosePrice)

	if len(b.Closes) < b.Period {
		return aReturn
	}

	upperBand, middleBand, lowerBand := b.GetTaLibBBands(b.Closes, b.Period)

	aReturn.BBandUpper = upperBand
	aReturn.BBandMiddle = middleBand
	aReturn.BBandLower = lowerBand

	return aReturn
}

func (b *BBands) GetTaLibBBands(closes []decimal.Decimal, period int) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	inReal := make([]float64, 0)
	for _, close := range closes {
		f, _ := close.Float64()
		inReal = append(inReal, f)
	}

	upperBands, middleBands, lowerBands := talib.BBands(inReal, period, 2, 2, 0)

	upperBand := upperBands[len(upperBands)-1]
	middleBand := middleBands[len(middleBands)-1]
	lowerBand := lowerBands[len(lowerBands)-1]

	return decimal.NewFromFloat(upperBand), decimal.NewFromFloat(middleBand), decimal.NewFromFloat(lowerBand)
}
