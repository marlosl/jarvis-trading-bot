package builders

import (
	"jarvis-trading-bot/structs"

	"github.com/shopspring/decimal"
)

type BotParametersBuilder struct {
	botParameters *structs.BotParameters // needs to be inited
}

func NewBotParametersBuilder() *BotParametersBuilder {
	return &BotParametersBuilder{&structs.BotParameters{}}
}

func (o *BotParametersBuilder) Build() *structs.BotParameters {
	return o.botParameters
}

func (o *BotParametersBuilder) SetUserId(v uint) *BotParametersBuilder {
	o.botParameters.UserId = v
	return o
}

func (o *BotParametersBuilder) SetBroker(v string) *BotParametersBuilder {
	o.botParameters.Broker = v
	return o
}

func (o *BotParametersBuilder) SetPercentageTax(v decimal.Decimal) *BotParametersBuilder {
	o.botParameters.PercentageTax = v
	return o
}

func (o *BotParametersBuilder) SetStopLossPercentage(v decimal.Decimal) *BotParametersBuilder {
	o.botParameters.StopLossPercentage = v
	return o
}

func (o *BotParametersBuilder) SetTrailingStopLoss(v bool) *BotParametersBuilder {
	o.botParameters.TrailingStopLoss = v
	return o
}

func (o *BotParametersBuilder) SetMinimumLimitPercentage(v decimal.Decimal) *BotParametersBuilder {
	o.botParameters.MinimumLimitPercentage = v
	return o
}

func (o *BotParametersBuilder) SetRsiPeriod(v int) *BotParametersBuilder {
	o.botParameters.RsiPeriod = v
	return o
}

func (o *BotParametersBuilder) SetRsiOverbought(v decimal.Decimal) *BotParametersBuilder {
	o.botParameters.RsiOverbought = v
	return o
}

func (o *BotParametersBuilder) SetRsiOversold(v decimal.Decimal) *BotParametersBuilder {
	o.botParameters.RsiOversold = v
	return o
}

func (o *BotParametersBuilder) SetMaxNumberNegotiations(v int) *BotParametersBuilder {
	o.botParameters.MaxNumberNegotiations = v
	return o
}

func (o *BotParametersBuilder) SetMinPeriodNextNegotiation(v int) *BotParametersBuilder {
	o.botParameters.MinPeriodNextNegotiation = v
	return o
}

func (o *BotParametersBuilder) SetSymbol(v string) *BotParametersBuilder {
	o.botParameters.Symbol = v
	return o
}

func (o *BotParametersBuilder) SetBuyingQty(v decimal.Decimal) *BotParametersBuilder {
	o.botParameters.BuyingQty = v
	return o
}

func (o *BotParametersBuilder) SetBuyingAsset(v string) *BotParametersBuilder {
	o.botParameters.BuyingAsset = v
	return o
}

func (o *BotParametersBuilder) SetSellingAsset(v string) *BotParametersBuilder {
	o.botParameters.SellingAsset = v
	return o
}

func (o *BotParametersBuilder) SetStreamSymbol(v string) *BotParametersBuilder {
	o.botParameters.StreamSymbol = v
	return o
}

func (o *BotParametersBuilder) SetStreamInterval(v int) *BotParametersBuilder {
	o.botParameters.StreamInterval = v
	return o
}

func (o *BotParametersBuilder) SetCreated(v *structs.Timestamp) *BotParametersBuilder {
	o.botParameters.Created = v
	return o
}

func (o *BotParametersBuilder) SetClosed(v *structs.Timestamp) *BotParametersBuilder {
	o.botParameters.Closed = v
	return o
}
