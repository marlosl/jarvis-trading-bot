package structs

import (
	"time"

	"github.com/shopspring/decimal"
)

type TradingStatusKey struct {
	UserId         uint
	BotParameterId uint
	InstanceId     uint
	Symbol         string
	Active         bool
}
type BotParams struct {
	Symbol                   string
	BuyingAsset              string
	SellingAsset             string
	BuyingQty                decimal.Decimal
	PercentageTax            decimal.Decimal
	StopLossPercentage       decimal.Decimal
	MinimumLimitPercentage   decimal.Decimal
	RsiPeriod                int
	RsiOverbought            decimal.Decimal
	RsiOversold              decimal.Decimal
	MaxNumberNegotiations    int
	MinPeriodNextNegotiation int
	TrailingStopLoss         bool
	StreamSymbol             string
	StreamInterval           int
}

type AnalysisReturn struct {
	Operation   string
	Price       decimal.Decimal
	PSI         decimal.Decimal
	BBandUpper  decimal.Decimal
	BBandMiddle decimal.Decimal
	BBandLower  decimal.Decimal
}

type CandlestickAnalysis struct {
	EventTime   time.Time       `json:"eventTime"`
	Symbol      string          `json:"symbol"`
	OpenPrice   decimal.Decimal `json:"openPrice"`
	ClosePrice  decimal.Decimal `json:"closePrice"`
	HighPrice   decimal.Decimal `json:"highPrice"`
	LowPrice    decimal.Decimal `json:"lowPrice"`
	Volume      decimal.Decimal `json:"volume"`
	Operation   string          `json:"operation"`
	PSI         decimal.Decimal `json:"psi"`
	BBandUpper  decimal.Decimal `json:"bbandUpper"`
	BBandMiddle decimal.Decimal `json:"bbandMiddle"`
	BBandLower  decimal.Decimal `json:"bbandLower"`
}
