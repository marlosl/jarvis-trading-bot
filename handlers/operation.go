package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"jarvis-trading-bot/services/types"

	"github.com/aws/aws-lambda-go/events"
)

func OperationHandler(_ context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		var transactionItem types.TransactionItem
		err := json.Unmarshal([]byte(message.Body), &transactionItem)
		if err != nil {
			fmt.Printf("Can't unmarshal body: %v", err)
			continue
		}
	}
	return nil
}
