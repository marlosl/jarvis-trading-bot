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
	v := GetDecimalConfig(consts.StopLossPerc)
	assert.NotEmpty(v)
}

func TestGetIntConfig(t *testing.T) {
	assert := assert.New(t)
	v := GetIntConfig(consts.PostgresPort)
	assert.NotEmpty(v)
}

func TestGetInt64Config(t *testing.T) {
	assert := assert.New(t)
	v := GetInt64Config(consts.PostgresPort)
	assert.NotEmpty(v)
}

func TestGetStringConfig(t *testing.T) {
	assert := assert.New(t)
	v := GetStringConfig(consts.BinanceUrl)
	assert.NotEmpty(v)
}

func TestGetBoleanConfig(t *testing.T) {
	assert := assert.New(t)
	v := GetBoleanConfig(consts.LogToFile)
	assert.NotEmpty(v)
}

func TestGetStringSliceConfig(t *testing.T) {
	assert := assert.New(t)
	v := GetStringSliceConfig(consts.Symbols)
	assert.NotEmpty(v)
}

func TestGetDecimalSliceConfig(t *testing.T) {
	assert := assert.New(t)
	v := GetDecimalSliceConfig(consts.BuyingQty)
	assert.NotEmpty(v)
}
