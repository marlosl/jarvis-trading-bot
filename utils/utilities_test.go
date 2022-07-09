package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func UtilitiesLoadFileEnvVars() {
	LoadEnvVarsFromFiles([]string{"../.env"})
}

func TestSPrintJson(t *testing.T) {
	obj := map[string]string{
		"test": "value",
	}

	assert := assert.New(t)
	v := SPrintJson(obj)
	assert.NotEmpty(v)
}

func TestCreateHash(t *testing.T) {
	assert := assert.New(t)
	v := CreateHash("secret", "data")
	assert.NotEmpty(v)
}

func TestConvertToTime(t *testing.T) {
	assert := assert.New(t)
	v := ConvertToTime(9797649039384)
	assert.NotEmpty(v)
}

func TestConvertToTimestamp(t *testing.T) {
	assert := assert.New(t)
	v := ConvertToTimestamp(time.Now())
	assert.NotEmpty(v)
}

func TestConvertInt64ToString(t *testing.T) {
	assert := assert.New(t)
	v := ConvertInt64ToString(1000)
	assert.Equal("1000", v)
}

func TestConvertStringToInt64(t *testing.T) {
	assert := assert.New(t)
	v := ConvertStringToInt64("1000")
	assert.NotEmpty(v)
}

func TestConvertStringToInt(t *testing.T) {
	assert := assert.New(t)
	v := ConvertStringToInt("10")
	assert.NotEmpty(v)
}

func TestGetDecimalValue(t *testing.T) {
	assert := assert.New(t)
	v := GetDecimalValue("10.00")
	assert.NotEmpty(v)
}
