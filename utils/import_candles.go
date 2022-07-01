package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/utils/log"

	"github.com/shopspring/decimal"
)

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}

func ImportDataFile() {
	importFile := os.Getenv("IMPORT_FILE")
	log.InfoLogger.Printf("opening file %s\n", importFile)

	f, err := os.Open(importFile)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v\n", err))
	}

	counter := 0
	r := bufio.NewReader(f)
	s, e := Readln(r)

	for e == nil {
		counter++
		stream := new(structs.Stream)
		if err := json.Unmarshal([]byte(s), &stream); err != nil {
			log.InfoLogger.Println("read:", err)
		}

		candle := ConvertStreamToCandle(stream)
		DB.Create(&candle)
		s, e = Readln(r)
	}

	log.InfoLogger.Printf("%d records imported\n", counter)
}

func ConvertStreamToCandle(stream *structs.Stream) *structs.Candlestick {
	candle := new(structs.Candlestick)
	candle.Symbol = stream.Data.Symbol
	candle.EventTime = ConvertToTime(stream.Data.EventTime)
	candle.NumberOfTrades = stream.Data.Kline.NumberOfTrades

	openPrice, _ := decimal.NewFromString(stream.Data.Kline.OpenPrice)
	candle.OpenPrice = openPrice

	lowPrice, _ := decimal.NewFromString(stream.Data.Kline.LowPrice)
	candle.LowPrice = lowPrice

	highPrice, _ := decimal.NewFromString(stream.Data.Kline.HighPrice)
	candle.HighPrice = highPrice

	closePrice, _ := decimal.NewFromString(stream.Data.Kline.ClosePrice)
	candle.ClosePrice = closePrice

	baseAssetVolume, _ := decimal.NewFromString(stream.Data.Kline.BaseAssetVolume)
	candle.BaseAssetVolume = baseAssetVolume

	return candle
}
