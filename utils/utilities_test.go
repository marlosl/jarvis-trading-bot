package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func UtilitiesLoadFileEnvVars() {
	InitConfig()
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

func TestGetDecimalValue(t *testing.T) {
	assert := assert.New(t)
	d := "10.00"
	v := GetDecimalValue(&d)
	assert.NotEmpty(v)
}
