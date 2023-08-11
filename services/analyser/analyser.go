package analyser

import (
	"fmt"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/services/types"
	"jarvis-trading-bot/utils"

	"github.com/shopspring/decimal"
)

type AnalyserService struct {
}

var ONE_HUNDRED = decimal.NewFromInt(100)

func NewAnalyserService() *AnalyserService {
	return &AnalyserService{}
}

func (a *AnalyserService) DoAnalysis(transactionItem *types.TransactionItem, signalItem *signal.SignalItem) *types.TransactionItem {
	if transactionItem == nil {
		fmt.Printf("Creating a new transaction for ticker %s\n", signalItem.Ticker)

		transactionItem = a.OpenTransactionIfByingSignal(signalItem)
		if transactionItem == nil {
			fmt.Println("Signal is not in buying state")
		}

		return transactionItem
	}

	fmt.Printf("Transaction: %v\n", transactionItem)

	if signalItem.Action == signal.PRICE_REQUEST && !a.CloseTransactionIfStopLossOrTakeProfit(signalItem, transactionItem) {
		return transactionItem
	}

	transactionItem.Signals = append(transactionItem.Signals, types.Signal{
		SignalTime:    utils.GetCurrentTime(),
		IndicatorName: signalItem.IndicatorName,
		Action:        signalItem.Action,
	})

	return a.CloseTransactionIfSellingSignal(signalItem, transactionItem)
}

func (a *AnalyserService) CloseTransactionIfStopLossOrTakeProfit(
	signalItem *signal.SignalItem,
	transactionItem *types.TransactionItem,
) bool {

	buyPrice := utils.GetDecimalValue(transactionItem.BuyPrice)
	fmt.Printf("Buy price: %v\n", buyPrice)

	lastPrice := utils.GetDecimalValue(signalItem.Close)
	fmt.Printf("Last price: %v\n", lastPrice)

	diff := lastPrice.Sub(buyPrice)
	fmt.Printf("Diff: %v\n", diff)

	if diff.Cmp(decimal.Zero) == -1 {
		fmt.Println("CloseTransactionIfStopLoss")
		return a.CloseTransactionIfStopLoss(diff, signalItem, transactionItem)
	}

	if diff.Cmp(decimal.Zero) == 1 {
		fmt.Println("CloseTransactionIfTakeProfit")
		return a.CloseTransactionIfTakeProfit(diff, signalItem, transactionItem)
	}
	return false
}

func (a *AnalyserService) CloseTransactionIfStopLoss(
	diff decimal.Decimal,
	signalItem *signal.SignalItem,
	transactionItem *types.TransactionItem,
) bool {

	buyPrice := utils.GetDecimalValue(transactionItem.BuyPrice)
	fmt.Printf("Buy price: %v\n", buyPrice)

	stopLoss := utils.GetDecimalValue(signalItem.StopLossPerc)
	fmt.Printf("Stop loss: %v\n", stopLoss)

	fmt.Printf("Percentage: %v\n", diff.Div(buyPrice).Abs().Mul(ONE_HUNDRED))

	if !stopLoss.IsZero() && diff.Div(buyPrice).Abs().Mul(ONE_HUNDRED).Cmp(stopLoss) >= 0 {
		fmt.Println("STOP LOSS state")
		signalItem.Action = signal.STOP_LOSS
		transactionItem.SellPrice = signalItem.Close
		transactionItem.Status = consts.STATUS_CLOSED
		return true
	}

	return false
}

func (a *AnalyserService) CloseTransactionIfTakeProfit(
	diff decimal.Decimal,
	signalItem *signal.SignalItem,
	transactionItem *types.TransactionItem,
) bool {

	buyPrice := utils.GetDecimalValue(transactionItem.BuyPrice)
	fmt.Printf("Buy price: %v\n", buyPrice)

	takeProfit := utils.GetDecimalValue(signalItem.TakeProfitPerc)
	fmt.Printf("Take profit: %v\n", takeProfit)

	fmt.Printf("Percentage: %v\n", diff.Div(buyPrice).Abs().Mul(ONE_HUNDRED))

	if !takeProfit.IsZero() && diff.Div(buyPrice).Abs().Mul(ONE_HUNDRED).Cmp(takeProfit) >= 0 {
		fmt.Println("TAKE PROFIT state")
		signalItem.Action = signal.TAKE_PROFIT
		transactionItem.SellPrice = signalItem.Close
		transactionItem.Status = consts.STATUS_CLOSED
		return true
	}
	return false
}

func (a *AnalyserService) OpenTransactionIfByingSignal(signalItem *signal.SignalItem) *types.TransactionItem {
	if signal.IsOpenAction(signalItem.Action) && signalItem.Close != nil {
		fmt.Println("Signal is in OPENING state")
		transactionItem := NewTransactionItem(signalItem)
		transactionItem.BuyPrice = signalItem.Close
		return transactionItem
	}
	return nil
}

func (a *AnalyserService) CloseTransactionIfSellingSignal(
	signalItem *signal.SignalItem,
	transactionItem *types.TransactionItem,
) *types.TransactionItem {

	if signal.IsCloseAction(signalItem.Action) && signalItem.Close != nil {
		fmt.Println("Signal is in CLOSING state")
		transactionItem.SellPrice = signalItem.Close
		transactionItem.Status = consts.STATUS_CLOSED
	}
	return transactionItem
}

func NewTransactionItem(signalItem *signal.SignalItem) *types.TransactionItem {
	signal := types.Signal{
		SignalTime:    utils.GetCurrentTime(),
		IndicatorName: signalItem.IndicatorName,
		Action:        signalItem.Action,
	}

	return &types.TransactionItem{
		Exchange: signalItem.Exchange,
		Ticker:   signalItem.Ticker,
		Interval: signalItem.Interval,
		Signals:  []types.Signal{signal},
		Status:   consts.STATUS_ACTIVE,
	}
}
