package signal

import (
	"time"
)

type SignalItem struct {
	Alert          string     `json:"alert" dynamodbav:"alert"`
	Exchange       *string    `json:"exchange,omitempty" dynamodbav:"exchange,omitempty"`
	Ticker         string     `json:"ticker" dynamodbav:"ticker"`
	Action         string     `json:"action" dynamodbav:"action"`
	IndicatorName  string     `json:"indicatorName" dynamodbav:"indicatorName"`
	Close          *string    `json:"close,omitempty" dynamodbav:"close,omitempty"`
	Open           *string    `json:"open,omitempty" dynamodbav:"open,omitempty"`
	High           *string    `json:"high,omitempty" dynamodbav:"high,omitempty"`
	Low            *string    `json:"low,omitempty" dynamodbav:"low,omitempty"`
	Time           *string    `json:"time,omitempty" dynamodbav:"time,omitempty"`
	Volume         *string    `json:"volume,omitempty" dynamodbav:"volume,omitempty"`
	Interval       string     `json:"interval" dynamodbav:"interval"`
	StopLossPerc   *string     `json:"stopLossPerc" dynamodbav:"stopLossPerc"`
	TakeProfitPerc *string     `json:"takeProfitPerc" dynamodbav:"takeProfitPerc"`
	TimeNow        *string    `json:"timenow,omitempty" dynamodbav:"timenow,omitempty"`
	Payload        *string    `json:"payload,omitempty" dynamodbav:"payload,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty" dynamodbav:"createdAt,omitempty"`
}

type SignalDbItem struct {
	PK         string     `json:"pk" dynamodbav:"PK"`
	SK         string     `json:"sk" dynamodbav:"SK"`
	Attributes SignalItem `json:"attributes" dynamodbav:"attributes"`
}
