package exchangeconfig

import (
	"time"
)

type ExchangeConfigItem struct {
	Exchange         string     `json:"exchange,omitempty" dynamodbav:"exchange"`
	Ticker           string     `json:"ticker" dynamodbav:"ticker"`
	Symbol           string     `json:"symbol" dynamodbav:"symbol"`
	BuyQty           string     `json:"buyQty" dynamodbav:"buyQty"`
	StopLossPerc     string     `json:"stopLossPerc" dynamodbav:"stopLossPerc"`
	TakeProfitPerc   string     `json:"takeProfitPerc" dynamodbav:"takeProfitPerc"`
	RealTransactions bool       `json:"realTransactions" dynamodbav:"realTransactions"`
	CreatedAt        *time.Time `json:"createdAt,omitempty" dynamodbav:"createdAt,omitempty"`
}

type ExchangeConfigDbItem struct {
	PK         string             `json:"pk" dynamodbav:"PK"`
	SK         string             `json:"sk" dynamodbav:"SK"`
	Attributes ExchangeConfigItem `json:"attributes" dynamodbav:"attributes"`
}
