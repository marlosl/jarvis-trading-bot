package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

func PrintJson(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func SPrintJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return string(b)
}

func GetDecimalValue(value *string) decimal.Decimal {
	if value == nil {
		return decimal.NewFromInt(0)
	}

	v, e := decimal.NewFromString(*value)
	if e != nil {
		return decimal.NewFromInt(0)
	}
	return v
}

func CreateStringError(s string) error {
	return errors.New(s)
}

func CreateHash(secret string, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func ConvertToTimestamp(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/1e6, 10)
}

func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func NewString(s string) *string {
	return &s
}

func NewTimeNow() *time.Time {
	t := GetCurrentTime()
	return &t
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
