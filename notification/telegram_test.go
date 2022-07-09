package notification

import (
	"testing"

	"jarvis-trading-bot/utils"

	"github.com/stretchr/testify/assert"
)

func TestTelegram(t *testing.T) {
	utils.LoadEnvVarsFromFiles([]string{"../.env"})
	v1 := SendJson("{\"test\":\"value\"}")
	assert.Equalf(t, true, v1, "Send JSON message")

	v2 := SendMessage("Test message", false)
	assert.Equalf(t, true, v2, "Send message")
}
