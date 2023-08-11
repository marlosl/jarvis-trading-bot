package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/services/transaction"
	"jarvis-trading-bot/services/types"
	"jarvis-trading-bot/utils"

	"github.com/aws/aws-lambda-go/events"
)

func SignalAnalyserHandler(_ context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		var sqsMessage types.SQSMessage
		err := json.Unmarshal([]byte(message.Body), &sqsMessage)
		if err != nil {
			fmt.Printf("Can't unmarshal body: %v", err)
			continue
		}

		fmt.Printf("Unmarshal sqsMessage: %s", utils.SPrintJson(sqsMessage))

		if sqsMessage.Message == "" {
			fmt.Printf("Message is empty")
			continue
		}

		var signal signal.SignalItem

		err = json.Unmarshal([]byte(sqsMessage.Message), &signal)
		if err != nil {
			fmt.Printf("Can't unmarshal sqsMessage: %v", err)
			continue
		}

		transactionService, err := transaction.NewTransactionService()
		if err != nil {
			fmt.Printf("Can't create transactionService: %v\n", err)
			return err
		}

		err = transactionService.DoAnalysis(&signal)
		if err != nil {
			fmt.Printf("Error while analysing item %s: %v\n", signal.Ticker, err)
		}
	}
	return nil
}
