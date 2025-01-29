package binance

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	Asset  string
	Free   decimal.Decimal `json:"free"`
	Locked decimal.Decimal `json:"locked"`
}

type AccountBalance struct {
	MakerCommission  int
	TakerCommission  int
	BuyerCommission  int
	SellerCommission int
	CanTrade         bool
	CanWithdraw      bool
	CanDeposit       bool
	UpdateTime       int64
	AccountType      string
	Balances         []Balance
	Permissions      []string
}

type Order struct {
	Symbol string
	Side   string
	Qty    decimal.Decimal
	Price  decimal.Decimal
}

type Kline struct {
	KlineOpenTime    string
	OpenPrice        string
	HighPrice        string
	LowPrice         string
	ClosePrice       string
	Volume           string
	KlineCloseTime   string
	QuoteAssetVolume string
}

type Fill struct {
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string
	TradeId         int
}

type OrderResponse struct {
	Symbol              string          `json:"symbol,omitempty"`
	OrderId             int64           `json:"orderId,omitempty"`
	OrderListId         int64           `json:"orderListId,omitempty"`
	ClientOrderId       string          `json:"clientOrderId,omitempty"`
	TransactTime        int64           `json:"transactionTime,omitempty"`
	Price               decimal.Decimal `json:"price,omitempty"`
	OrigQty             decimal.Decimal `json:"origQty,omitempty"`
	ExecutedQty         decimal.Decimal `json:"executedQty,omitempty"`
	CummulativeQuoteQty decimal.Decimal `json:"cummulativeQuoteQty,omitempty"`
	Status              string          `json:"status,omitempty"`
	TimeInForce         string          `json:"timeInForce,omitempty"`
	Type                string          `json:"type,omitempty"`
	Side                string          `json:"side,omitempty"`
	Fills               []Fill          `json:"fills,omitempty"`
}

type TradeResponse struct {
	Symbol          string          `json:"symbol,omitempty"`
	Id              int64           `json:"id,omitempty"`
	OrderId         int64           `json:"orderId,omitempty"`
	OrderListId     int64           `json:"orderListId,omitempty"`
	Price           decimal.Decimal `json:"price,omitempty"`
	Qty             decimal.Decimal `json:"qty,omitempty"`
	QuoteQty        decimal.Decimal `json:"quoteQty,omitempty"`
	Commission      decimal.Decimal `json:"commission,omitempty"`
	CommissionAsset string          `json:"commissionAsset,omitempty"`
	Time            int64           `json:"time,omitempty"`
	IsBuyer         bool            `json:"isBuyer,omitempty"`
	IsMaker         bool            `json:"isMaker,omitempty"`
	IsBestMatch     bool            `json:"isBestMatch,omitempty"`
}
