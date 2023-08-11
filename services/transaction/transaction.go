package transaction

import (
	"fmt"
	"os"

	"jarvis-trading-bot/clients/topic"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/analyser"
	"jarvis-trading-bot/services/exchangeconfig"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/services/tickerconfig"
	"jarvis-trading-bot/services/types"
)

type TransactionService struct {
	repository *TransactionRepository
	analyser   *analyser.AnalyserService
	snsClient  *topic.SNSClient
}

func NewTransactionService() (*TransactionService, error) {
	repository, err := NewTransactionRepository(nil)
	if err != nil {
		fmt.Printf("Can't create dbClient: %v", err)
		return nil, err
	}

	operationTopic := os.Getenv(consts.SNSOperationTopic)
	snsClient, err := topic.NewSNSClient(&operationTopic)
	if err != nil {
		fmt.Printf("Can't initialize the topic %s: %v\n", operationTopic, err)
		return nil, err
	}

	return &TransactionService{
		repository: repository,
		analyser:   analyser.NewAnalyserService(),
		snsClient:  snsClient,
	}, nil
}

func (t *TransactionService) DoAnalysis(signalItem *signal.SignalItem) error {
	items, err := t.repository.GetItemsByStatus(signalItem.Ticker, consts.STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("Can't get items: %v", err)
		return err
	}

	if len(items) == 0 {
		fmt.Printf("No active transactions for ticker %s\n", signalItem.Ticker)

		item := t.analyser.DoAnalysis(nil, signalItem)
		if item == nil {
			fmt.Println("item is nil")
			return nil
		}

		err := t.repository.SaveItem(item)
		if err != nil {
			fmt.Printf("Can't save new item: %v\n", err)
		}

		err = tickerconfig.CreateTickerConfigIfNotExists(signalItem.Ticker)
		if err != nil {
			fmt.Printf("Can't save ticker config: %v\n", err)
		}

		t.UpdateTransactionItem(item)

		if item.RealTransactions && item.Status == consts.STATUS_ACTIVE {
			err = t.snsClient.SendMsg(&signalItem)
			if err != nil {
				fmt.Printf("Can't send the message to the queue: %v\n", err)
			}
		}
		return err
	}

	for _, item := range items {
		fmt.Printf("Transaction: %v", item)

		item := t.analyser.DoAnalysis(&item, signalItem)

		err := t.repository.SaveItem(item)
		if err != nil {
			fmt.Printf("Can't save item: %v", err)
			return err
		}

		if item.RealTransactions && item.Status == consts.STATUS_CLOSED {
			err = t.snsClient.SendMsg(&signalItem)
			if err != nil {
				fmt.Printf("Can't send the message to the queue: %v\n", err)
			}
		}
	}

	return nil
}

func (t *TransactionService) GetActiveItems() ([]types.TransactionItem, error) {
	return t.repository.GetActiveItems()
}

func (t *TransactionService) UpdateTransactionItem(item *types.TransactionItem) {
	config, err := exchangeconfig.GetExchangeConfig(item.Ticker, *item.Exchange)
	if err != nil {
		fmt.Printf("Can't get exchange config: %v\n", err)
	}

	if config != nil {
		item.RealTransactions = config.RealTransactions
		item.BuyQty = config.BuyQty
		item.Symbol = &config.Symbol
	}
}
