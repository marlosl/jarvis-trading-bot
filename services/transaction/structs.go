package transaction

import (
	"jarvis-trading-bot/services/types"
)

type TransactionDbItem struct {
	PK         string                `json:"pk" dynamodbav:"PK"`
	SK         string                `json:"sk" dynamodbav:"SK"`
	GS1PK      string                `json:"gs1pk" dynamodbav:"GS1PK"`
	GS1SK      string                `json:"gs1sk" dynamodbav:"GS1SK"`
	GS2PK      string                `json:"gs2pk" dynamodbav:"GS2PK"`
	GS2SK      string                `json:"gs2sk" dynamodbav:"GS2SK"`
	Attributes types.TransactionItem `json:"attributes" dynamodbav:"attributes"`
}
