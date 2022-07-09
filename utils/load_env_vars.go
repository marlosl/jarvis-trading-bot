package utils

import (
	"os"
	"strconv"
	"strings"

	"jarvis-trading-bot/utils/log"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func LoadEnvVars() {
	files := []string{".parameters", ".env"}
	LoadEnvVarsFromFiles(files)
}

func LoadEnvVarsFromFiles(files []string) {
	validList := make([]string, 0)
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			continue
		}
		validList = append(validList, f)
	}

	if err := godotenv.Load(validList...); err != nil {
		log.ErrorLogger.Fatal("Error loading .env file")
	}
}

func GetDecimalConfig(confName string) decimal.Decimal {
	s := os.Getenv(confName)
	v, e := decimal.NewFromString(s)
	if e != nil {
		return decimal.NewFromInt(0)
	}
	return v
}

func GetIntConfig(confName string) int {
	s := os.Getenv(confName)
	v, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return v
}

func GetInt64Config(confName string) int64 {
	return int64(GetIntConfig(confName))
}

func GetStringConfig(confName string) string {
	s := os.Getenv(confName)
	return s
}

func GetBoleanConfig(confName string) bool {
	s := os.Getenv(confName)
	return s == "true"
}

func GetStringSliceConfig(confName string) []string {
	s := os.Getenv(confName)
	return strings.Split(s, ",")
}

func GetDecimalSliceConfig(confName string) []decimal.Decimal {
	l := make([]decimal.Decimal, 0)
	s := GetStringSliceConfig(confName)
	for _, v := range s {
		d, e := decimal.NewFromString(v)
		if e != nil {
			d = decimal.NewFromInt(0)
		}
		l = append(l, d)
	}
	return l
}
