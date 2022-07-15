package builders

import (
	"jarvis-trading-bot/structs"

	"github.com/shopspring/decimal"
)

type BotParamsBuilder struct {
	botParams *structs.BotParams // needs to be inited
}

func NewBotParamsBuilder() *BotParamsBuilder {
	return &BotParamsBuilder{&structs.BotParams{}}
}

func (o *BotParamsBuilder) Build() *structs.BotParams {
	return o.botParams
}

func (o *BotParamsBuilder) SetPercentageTax(v decimal.Decimal) *BotParamsBuilder {
	o.botParams.PercentageTax = v
	return o
}

func (o *BotParamsBuilder) SetStopLossPercentage(v decimal.Decimal) *BotParamsBuilder {
	o.botParams.StopLossPercentage = v
	return o
}

func (o *BotParamsBuilder) SetTrailingStopLoss(v bool) *BotParamsBuilder {
	o.botParams.TrailingStopLoss = v
	return o
}

func (o *BotParamsBuilder) SetMinimumLimitPercentage(v decimal.Decimal) *BotParamsBuilder {
	o.botParams.MinimumLimitPercentage = v
	return o
}

func (o *BotParamsBuilder) SetRsiPeriod(v int) *BotParamsBuilder {
	o.botParams.RsiPeriod = v
	return o
}

func (o *BotParamsBuilder) SetRsiOverbought(v decimal.Decimal) *BotParamsBuilder {
	o.botParams.RsiOverbought = v
	return o
}

func (o *BotParamsBuilder) SetRsiOversold(v decimal.Decimal) *BotParamsBuilder {
	o.botParams.RsiOversold = v
	return o
}

func (o *BotParamsBuilder) SetMaxNumberNegotiations(v int) *BotParamsBuilder {
	o.botParams.MaxNumberNegotiations = v
	return o
}

func (o *BotParamsBuilder) SetMinPeriodNextNegotiation(v int) *BotParamsBuilder {
	o.botParams.MinPeriodNextNegotiation = v
	return o
}

func (o *BotParamsBuilder) SetSymbol(v string) *BotParamsBuilder {
	o.botParams.Symbol = v
	return o
}

func (o *BotParamsBuilder) SetBuyingQty(v decimal.Decimal) *BotParamsBuilder {
	o.botParams.BuyingQty = v
	return o
}

func (o *BotParamsBuilder) SetBuyingAsset(v string) *BotParamsBuilder {
	o.botParams.BuyingAsset = v
	return o
}

func (o *BotParamsBuilder) SetSellingAsset(v string) *BotParamsBuilder {
	o.botParams.SellingAsset = v
	return o
}

func (o *BotParamsBuilder) SetStreamSymbol(v string) *BotParamsBuilder {
	o.botParams.StreamSymbol = v
	return o
}

func (o *BotParamsBuilder) SetStreamInterval(v int) *BotParamsBuilder {
	o.botParams.StreamInterval = v
	return o
}
