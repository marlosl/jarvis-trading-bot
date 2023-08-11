package types

import (
	"time"
)

type Signal struct {
	SignalTime    time.Time `json:"signalTime" dynamodbav:"signalTime"`
	IndicatorName string    `json:"indicatorName" dynamodbav:"indicatorName"`
	Action        string    `json:"action" dynamodbav:"action"`
}

type TransactionItem struct {
	Uuid             string     `json:"uuid" dynamodbav:"uuid"`
	Exchange         *string    `json:"exchange,omitempty" dynamodbav:"exchange,omitempty"`
	Ticker           string     `json:"ticker" dynamodbav:"ticker"`
	Symbol           *string    `json:"symbol,omitempty" dynamodbav:"symbol,omitempty"`
	Signals          []Signal   `json:"signals" dynamodbav:"signals"`
	BuyPrice         *string    `json:"buyPrice,omitempty" dynamodbav:"buyPrice,omitempty"`
	SellPrice        *string    `json:"sellPrice,omitempty" dynamodbav:"sellPrice,omitempty"`
	Status           string     `json:"status" dynamodbav:"status"`
	Interval         string     `json:"interval" dynamodbav:"interval"`
	RealTransactions bool       `json:"realTransactions" dynamodbav:"realTransactions"`
	BuyQty           string     `json:"buyQty" dynamodbav:"buyQty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty" dynamodbav:"createdAt,omitempty"`
}

type TransactionUpdateItem struct {
	Signals   []Signal `json:":signals" dynamodbav:":signals"`
	BuyPrice  *string  `json:":buyPrice,omitempty" dynamodbav:":buyPrice,omitempty"`
	SellPrice *string  `json:":sellPrice,omitempty" dynamodbav:":sellPrice,omitempty"`
	Status    string   `json:":status" dynamodbav:":status,updated"`
}

type SQSMessage struct {
	Message string `json:"Message"`
}
