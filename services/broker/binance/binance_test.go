package binance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinanceGetCandlestick(t *testing.T) {
	assert := assert.New(t)

	api := &BinanceApi{
		BaseUrl: "https://data.binance.com",
	}

	candle, err := api.GetCandlestick("BTCUSDT")
	if err != nil {
		t.Log(err)
	}
	assert.Nil(err)
	assert.NotEmpty(candle.OpenPrice)
}
