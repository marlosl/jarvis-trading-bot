package builders

import (
	"jarvis-trading-bot/structs"
	"time"

	"github.com/shopspring/decimal"
)

type TradingStatusBuilder struct {
	tradingStatus *structs.TradingStatus // needs to be inited
}

func NewTradingStatusBuilder() *TradingStatusBuilder {
	return &TradingStatusBuilder{&structs.TradingStatus{}}
}

func (o *TradingStatusBuilder) Build() *structs.TradingStatus {
	return o.tradingStatus
}

func (o *TradingStatusBuilder) SetUserId(v uint) *TradingStatusBuilder {
	o.tradingStatus.UserId = v
	return o
}

func (o *TradingStatusBuilder) SetBotParameterId(v uint) *TradingStatusBuilder {
	o.tradingStatus.BotParameterId = v
	return o
}

func (o *TradingStatusBuilder) SetInstanceId(v uint) *TradingStatusBuilder {
	o.tradingStatus.InstanceId = v
	return o
}

func (o *TradingStatusBuilder) SetSymbol(v string) *TradingStatusBuilder {
	o.tradingStatus.Symbol = v
	return o
}

func (o *TradingStatusBuilder) SetLastStatus(v int) *TradingStatusBuilder {
	o.tradingStatus.LastStatus = v
	return o
}

func (o *TradingStatusBuilder) SetBuyAmount(v decimal.Decimal) *TradingStatusBuilder {
	o.tradingStatus.BuyAmount = v
	return o
}

func (o *TradingStatusBuilder) SetInitialBuyAmount(v decimal.Decimal) *TradingStatusBuilder {
	o.tradingStatus.InitialBuyAmount = v
	return o
}

func (o *TradingStatusBuilder) SetProfitAmount(v decimal.Decimal) *TradingStatusBuilder {
	o.tradingStatus.ProfitAmount = v
	return o
}

func (o *TradingStatusBuilder) SetBuyTax(v decimal.Decimal) *TradingStatusBuilder {
	o.tradingStatus.BuyTax = v
	return o
}

func (o *TradingStatusBuilder) SetSellTax(v decimal.Decimal) *TradingStatusBuilder {
	o.tradingStatus.SellTax = v
	return o
}

func (o *TradingStatusBuilder) SetTotalProfit(v decimal.Decimal) *TradingStatusBuilder {
	o.tradingStatus.TotalProfit = v
	return o
}

func (o *TradingStatusBuilder) SetTotalTaxes(v decimal.Decimal) *TradingStatusBuilder {
	o.tradingStatus.TotalTaxes = v
	return o
}

func (o *TradingStatusBuilder) SetLastEvent(v time.Time) *TradingStatusBuilder {
	o.tradingStatus.LastEvent = v
	return o
}

func (o *TradingStatusBuilder) SetLastOperationTime(v time.Time) *TradingStatusBuilder {
	o.tradingStatus.LastOperationTime = v
	return o
}

func (o *TradingStatusBuilder) SetSimulation(v bool) *TradingStatusBuilder {
	o.tradingStatus.Simulation = v
	return o
}

func (o *TradingStatusBuilder) SetCreatedAt(v time.Time) *TradingStatusBuilder {
	o.tradingStatus.CreatedAt = v
	return o
}

func (o *TradingStatusBuilder) SetUpdatedAt(v time.Time) *TradingStatusBuilder {
	o.tradingStatus.UpdatedAt = v
	return o
}
