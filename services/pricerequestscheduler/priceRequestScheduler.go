package pricerequestscheduler

import (
	"fmt"
	"os"

	"jarvis-trading-bot/clients/topic"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/services/transaction"
	"jarvis-trading-bot/utils"
)

func CreatePriceRequests() error {

	transactionService, err := transaction.NewTransactionService()
	if err != nil {
		fmt.Printf("Can't create transactionService: %v", err)
		return err
	}

	priceRequestTopic := os.Getenv(consts.SNSPriceRequestTopic)
	snsClient, err := topic.NewSNSClient(&priceRequestTopic)
	if err != nil {
		fmt.Printf("Can't initialize the topic %s: %v\n", priceRequestTopic, err)
		return err
	}

	items, err := transactionService.GetActiveItems()
	if err != nil {
		fmt.Printf("Can't get items: %v", err)
		return err
	}

	for _, item := range items {
		signalItem := &signal.SignalItem{
			Alert:         "",
			Exchange:      item.Exchange,
			Ticker:        item.Ticker,
			Action:        signal.PRICE_REQUEST,
			IndicatorName: "",
			Interval:      item.Interval,
		}

		fmt.Println("Signal: ", utils.SPrintJson(signalItem))

		err = snsClient.SendMsg(&signalItem)
		if err != nil {
			fmt.Printf("Can't send the message to the queue: %v\n", err)
		}
	}
	return nil
}
