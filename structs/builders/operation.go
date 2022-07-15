package builders

import (
	"jarvis-trading-bot/structs"
	"time"

	"github.com/shopspring/decimal"
)

type OperationBuilder struct {
	operation *structs.Operation // needs to be inited
}

func NewOperationBuilder() *OperationBuilder {
	return &OperationBuilder{&structs.Operation{}}
}

func (o *OperationBuilder) Build() *structs.Operation {
	return o.operation
}

func (o *OperationBuilder) SetBotParameterId(v uint) *OperationBuilder {
	o.operation.BotParameterId = v
	return o
}

func (o *OperationBuilder) SetInstanceId(v uint) *OperationBuilder {
	o.operation.InstanceId = v
	return o
}

func (o *OperationBuilder) SetSymbol(v string) *OperationBuilder {
	o.operation.Symbol = v
	return o
}

func (o *OperationBuilder) SetOperation(v string) *OperationBuilder {
	o.operation.Operation = v
	return o
}

func (o *OperationBuilder) SetBaseAsset(v string) *OperationBuilder {
	o.operation.BaseAsset = v
	return o
}

func (o *OperationBuilder) SetBasePrice(v decimal.Decimal) *OperationBuilder {
	o.operation.BasePrice = v
	return o
}

func (o *OperationBuilder) SetOrderId(v int64) *OperationBuilder {
	o.operation.OrderId = v
	return o
}

func (o *OperationBuilder) SetOrigQty(v decimal.Decimal) *OperationBuilder {
	o.operation.OrigQty = v
	return o
}

func (o *OperationBuilder) SetExecutedQty(v decimal.Decimal) *OperationBuilder {
	o.operation.ExecutedQty = v
	return o
}

func (o *OperationBuilder) SetCummulativeQuoteQty(v decimal.Decimal) *OperationBuilder {
	o.operation.CummulativeQuoteQty = v
	return o
}

func (o *OperationBuilder) SetCommissionBase(v decimal.Decimal) *OperationBuilder {
	o.operation.CommissionBase = v
	return o
}

func (o *OperationBuilder) SetCommission(v decimal.Decimal) *OperationBuilder {
	o.operation.Commission = v
	return o
}

func (o *OperationBuilder) SetCommissionAsset(v string) *OperationBuilder {
	o.operation.CommissionAsset = v
	return o
}

func (o *OperationBuilder) SetType(v string) *OperationBuilder {
	o.operation.Type = v
	return o
}

func (o *OperationBuilder) SetStatus(v string) *OperationBuilder {
	o.operation.Status = v
	return o
}

func (o *OperationBuilder) SetOpened(v time.Time) *OperationBuilder {
	o.operation.Opened = v
	return o
}

func (o *OperationBuilder) SetTransactTime(v time.Time) *OperationBuilder {
	o.operation.TransactTime = v
	return o
}

func (o *OperationBuilder) SetFinished(v bool) *OperationBuilder {
	o.operation.Finished = v
	return o
}

func (o *OperationBuilder) SetSimulation(v bool) *OperationBuilder {
	o.operation.Simulation = v
	return o
}
