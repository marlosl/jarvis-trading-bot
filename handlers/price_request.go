package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"jarvis-trading-bot/clients/topic"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/pricerequest"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/services/types"
	"jarvis-trading-bot/utils"

	"github.com/aws/aws-lambda-go/events"
)

func PriceRequestHandler(_ context.Context, sqsEvent events.SQSEvent) error {
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

		err = pricerequest.DoPriceRequest(&signal)
		if err != nil {
			fmt.Printf("Can't do price request: %v", err)
			continue
		}

		signalTopic := os.Getenv(consts.SNSSignalsTopic)
		snsClient, err := topic.NewSNSClient(&signalTopic)
		if err != nil {
			fmt.Printf("Can't SNS Client: %v", err)
			continue
		}

		err = snsClient.SendMsg(&signal)
		if err != nil {
			fmt.Printf("Can't send SNS message: %v", err)
			continue
		}
	}
	return nil
}
