package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/transaction"
	"jarvis-trading-bot/utils"

	"github.com/kataras/tablewriter"
	"github.com/lensesio/tableprinter"
	"github.com/shopspring/decimal"
)

type TransactionItem struct {
	Ticker    string `header:"Ticker"`
	BuyPrice  string `header:"BuyPrice"`
	SellPrice string `header:"SellPrice"`
	Profit    string `header:"Profit"`
	Status    string `header:"Status"`
	Action    string `header:"Action"`
	CreatedAt string `header:"CreatedAt"`
}

func GetTransactions(tickers []string) {
	if len(tickers) == 0 {
		tickers = append(tickers, "BTCBRL", "BTCUSD")
	}

	results := getResultsAndSort(tickers...)

	printer := tableprinter.New(os.Stdout)

	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.HeaderBgColor = tablewriter.BgBlackColor
	printer.HeaderFgColor = tablewriter.FgGreenColor

	printer.Print(results)
}

func GenerateFile(tickers []string, filename string) {
	results := getResultsAndSort(tickers...)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to open file: %v", err)
		return
	}

	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Comma = ';'

	header := GetHeaders(TransactionItem{})
	if err := w.Write(header); err != nil {
		fmt.Printf("Error writing header to file: %v", err)
		return
	}

	for _, r := range results {
		record := GetValues(r)
		if err := w.Write(record); err != nil {
			fmt.Printf("Error writing record to file: %v", err)
		}
	}
}

func getResultsAndSort(tickers ...string) (results []TransactionItem) {
	repository, err := transaction.NewTransactionRepository(nil)
	if err != nil {
		return
	}

	results = make([]TransactionItem, 0)
	for _, ticker := range tickers {
		results = getResults(repository, ticker, consts.STATUS_ACTIVE, &results)
		results = getResults(repository, ticker, consts.STATUS_CLOSED, &results)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[j].Ticker > results[i].Ticker && results[j].CreatedAt > results[i].CreatedAt
	})

	return results
}

func getResults(repository *transaction.TransactionRepository, ticker string, status string, r *[]TransactionItem) []TransactionItem {
	var results []TransactionItem
	if r == nil {
		results = make([]TransactionItem, 0)
	} else {
		results = *r
	}

	items, err :=
		repository.GetItemsByStatus(ticker, status)
	if err != nil {
		fmt.Printf("Can't get items: %v", err)
		return *r
	}

	for _, item := range items {
	    lastSignal := item.Signals[len(item.Signals)-1]
		buyPrice := utils.GetDecimalValue(item.BuyPrice)
		sellPrice := utils.GetDecimalValue(item.SellPrice)
		profit := decimal.Zero

		if status == consts.STATUS_CLOSED {
			profit = sellPrice.Sub(buyPrice)
		}

		result := &TransactionItem{
			Ticker:    item.Ticker,
			BuyPrice:  buyPrice.String(),
			SellPrice: sellPrice.String(),
			Profit:    profit.String(),
			Status:    item.Status,
			Action:    lastSignal.Action,
			CreatedAt: utils.FormatDateTime(*item.CreatedAt),
		}
		results = append(results, *result)
	}

	return results
}

func GetHeaders(s interface{}) []string {
	values := []string{}

	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return values
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("header"), ",")

		fmt.Println("Name:", f.Name)
		if len(v) > 0 {
			values = append(values, v[0])
		}
	}
	return values
}

func GetValues(s interface{}) []string {
	values := []string{}

	if reflect.TypeOf(s).Kind() != reflect.Struct {
		return values
	}

	e := reflect.ValueOf(&s).Elem()
	for i := 0; i < e.NumField(); i++ {
		v := e.Field(i).Interface()
		values = append(values, v.(string))
	}
	return values
}
