package pricerequest

import (
	"errors"
	"fmt"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/broker/binance"
	"jarvis-trading-bot/services/types"
	"jarvis-trading-bot/utils"
)

func DoOperation(item *types.TransactionItem) error {
	if item == nil {
		fmt.Println("item is nil")
		return errors.New("item is nil")
	}

	if !item.RealTransactions {
		fmt.Printf("This item %s is not for real transactions.\n", item.Ticker)
		return errors.New("item is not for real transactions")
	}

	var order *binance.OrderResponse
	var err error

	api := binance.NewBinanceApi()

	if item.Status == consts.STATUS_ACTIVE {
		order, err = api.Buy(*item.Symbol, utils.GetDecimalValue(&item.BuyQty), utils.GetDecimalValue(item.BuyPrice))

		if err != nil {
			fmt.Printf("Can't buy: %v", err)
			return err
		}

		*item.BuyPrice = order.Price.String()
		fmt.Printf("Successfully bought %s for %s\n", *item.Symbol, *item.BuyPrice)
	}

	if item.Status == consts.STATUS_CLOSED {
		order, err = api.Sell(*item.Symbol, utils.GetDecimalValue(&item.BuyQty), utils.GetDecimalValue(item.SellPrice))

		if err != nil {
			fmt.Printf("Can't sell: %v", err)
			return err
		}

		*item.SellPrice = order.Price.String()
		fmt.Printf("Successfully sell %s for %s\n", *item.Symbol, *item.BuyPrice)
	}
	return nil
}
