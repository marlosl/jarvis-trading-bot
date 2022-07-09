package utils

import (
	"testing"

	"jarvis-trading-bot/consts"

	"github.com/stretchr/testify/assert"
)

func LoadFileEnvVars() {
	LoadEnvVarsFromFiles([]string{"../.env"})
}
func TestGetDecimalConfig(t *testing.T) {
	assert := assert.New(t)
	GetDecimalConfig(consts.StopLossPerc)
	assert.NotNil(DB)
}

func TestGetIntConfig(t *testing.T) {
	assert := assert.New(t)

	GetIntConfig(consts.PostgresPort)
	assert.NotNil(DB)
}

func TestGetInt64Config(t *testing.T) {
	assert := assert.New(t)

	GetInt64Config(consts.PostgresPort)
	assert.NotEmpty(DB)
}

func TestGetStringConfig(t *testing.T) {
	assert := assert.New(t)

	GetStringConfig(consts.BinanceUrl)
	assert.NotNil(DB)
}

func TestGetBoleanConfig(t *testing.T) {
	assert := assert.New(t)

	GetBoleanConfig(consts.LogToFile)
	assert.NotNil(DB)
}

func TestGetStringSliceConfig(t *testing.T) {
	assert := assert.New(t)

	GetStringSliceConfig(consts.Symbols)
	assert.NotNil(DB)
}

func TestGetDecimalSliceConfig(t *testing.T) {
	assert := assert.New(t)

	GetDecimalSliceConfig(consts.BuyingQty)
	assert.NotNil(DB)
}
