package broker

import (
	"testing"

	"jarvis-trading-bot/utils"

	"github.com/stretchr/testify/assert"
)

func TestBinance(t *testing.T) {
	utils.LoadEnvVarsFromFiles([]string{"../.env"})

	b := new(BinanceApi)
	b.Init()

	b1 := b.GetAssetBalance("BTC")
	assert.NotNilf(t, b1, "Balance is not nil")

	b2, err := b.GetBalance()
	assert.Nilf(t, err, "Account Balance has no errors")
	assert.NotNilf(t, b2, "Account Balance is not nil")
}
