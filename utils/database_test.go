package utils

import (
	"os"
	"testing"

	"jarvis-trading-bot/consts"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	assert := assert.New(t)

	LoadEnvVarsFromFiles([]string{"../.env"})
	os.Setenv(consts.SQLiteDatabase, "../candles.db")
	InitDatabase()
	assert.NotNil(DB)
}
