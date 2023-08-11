package tickerconfig

import (
	"time"
)

type TickerConfigItem struct {
	Ticker    string     `json:"ticker" dynamodbav:"ticker"`
	CreatedAt *time.Time `json:"createdAt,omitempty" dynamodbav:"createdAt,omitempty"`
}

type TickerConfigDbItem struct {
	PK         string           `json:"pk" dynamodbav:"PK"`
	SK         string           `json:"sk" dynamodbav:"SK"`
	Attributes TickerConfigItem `json:"attributes" dynamodbav:"attributes"`
}
