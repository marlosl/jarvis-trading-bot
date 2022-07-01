package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"jarvis-trading-bot/utils/log"

	"github.com/shopspring/decimal"
)

func PrintJson(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	log.InfoLogger.Println(string(b))
}

func SPrintJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return string(b)
}

func CreateHash(secret string, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func ConvertToTime(ts int64) time.Time {
	f := ts / 1000
	i := int64(f)
	return time.Unix(i, 0)
}

func ConvertToTimestamp(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/1e6, 10)
}

func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ConvertStringToInt64(s string) int64 {
	if i, e := strconv.ParseInt(s, 10, 64); e == nil {
		return i
	}
	return 0
}

func ConvertStringToInt(s string) int {
	if i, e := strconv.Atoi(s); e == nil {
		return i
	}
	return 0
}

func GetDecimalValue(value string) decimal.Decimal {
	v, e := decimal.NewFromString(value)
	if e != nil {
		return decimal.NewFromInt(0)
	}
	return v
}
