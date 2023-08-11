package signal

import (
	"strings"

	"jarvis-trading-bot/utils"
)

const (
	ALERT    = "alert"
	EXCHANGE = "exchange"
	TICKER   = "ticker"
	CLOSE    = "close"
	OPEN     = "open"
	HIGH     = "high"
	LOW      = "low"
	TIME     = "time"
	INTERVAL = "interval"
	VOLUME   = "volume"
	TIMENOW  = "timenow"

	BUY           = "BUY"
	SELL          = "SELL"
	POSITION      = "POSITION"
	TOP           = "TOP"
	BOTTOM        = "BOTTOM"
	REVERSAL      = "REVERSAL"
	TREND_CHANGED = "TRENDCHANGED"

	STOP_LOSS     = "STOPLOSS"
	TAKE_PROFIT   = "TAKEPROFIT"
	PRICE_REQUEST = "PRICEREQUEST"

	ALGOD       = "ALGOD"
	MK_PRO      = "MK_PRO"
	MK_PRO_1    = "MK-PRO"
	MK_PRO_2    = "MK PRO"
	MK_PRO_3    = "MKPRO"
	ALGO_PRO_V1 = "ALGOPROV1"
	ALGO_PRO_V2 = "ALGOPROV2"
	ALGO_PRO_V3 = "ALGOPROV3"
	ALGO_SWING  = "ALGOSWING"

	INTERVAL_1M  = "1M"
	INTERVAL_5M  = "5M"
	INTERVAL_10M = "10M"
	INTERVAL_30M = "30M"
	INTERVAL_1H  = "1H"
	INTERVAL_1D  = "1D"

	NONE = "NONE"
)

func ConvertTextToSignal(text string) (*SignalItem, error) {
	mapSignal := map[string]string{}

	chunks := strings.Split(text, ",")
	for _, chunk := range chunks {
		if chunk == "" {
			continue
		}

		parts := strings.Split(chunk, ": ")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		mapSignal[key] = value
	}

	if len(mapSignal) == 0 {
		return nil, utils.CreateStringError("Can't convert text to signal")
	}

	alert := mapSignal[ALERT]
	exchange := mapSignal[EXCHANGE]
	close := mapSignal[CLOSE]
	open := mapSignal[OPEN]
	high := mapSignal[HIGH]
	low := mapSignal[LOW]
	time := mapSignal[TIME]
	volume := mapSignal[VOLUME]
	timeNow := mapSignal[TIMENOW]

	return &SignalItem{
		Alert:         alert,
		Exchange:      &exchange,
		Ticker:        mapSignal[TICKER],
		Action:        GetActionFromAlert(alert),
		IndicatorName: GetIndicatorNameFromAlert(alert),
		Close:         &close,
		Open:          &open,
		High:          &high,
		Low:           &low,
		Time:          &time,
		Volume:        &volume,
		Interval:      GetIntervalFromAlert(alert),
		TimeNow:       &timeNow,
		Payload:       &text,
	}, nil
}

func ConvertAlertToSignal(alert string, text string) (*SignalItem, error) {
	ticker := ""

	chunks := strings.Split(alert, "-")
	if len(chunks) > 1 {
		ticker = strings.TrimSpace(chunks[len(chunks)-1])
	}

	if ticker == "" {
		return nil, utils.CreateStringError("Can't convert alert to signal")
	}

	return &SignalItem{
		Alert:         alert,
		Exchange:      nil,
		Ticker:        ticker,
		Action:        GetActionFromAlert(text),
		IndicatorName: GetIndicatorNameFromAlert(alert),
		Close:         nil,
		Open:          nil,
		High:          nil,
		Low:           nil,
		Time:          nil,
		Volume:        nil,
		Interval:      GetIntervalFromAlert(alert),
		TimeNow:       nil,
		Payload:       &text,
	}, nil
}

func GetActionFromAlert(alert string) string {
	s := strings.ToUpper(alert)
	switch {
	case strings.Contains(s, BUY):
		return BUY
	case strings.Contains(s, SELL):
		return SELL
	case strings.Contains(s, POSITION):
		return POSITION
	case strings.Contains(s, TOP):
		return TOP
	case strings.Contains(s, BOTTOM):
		return BOTTOM
	case strings.Contains(s, REVERSAL):
		return REVERSAL
	case strings.Contains(s, TREND_CHANGED):
		return TREND_CHANGED
	}

	return NONE
}

func GetIndicatorNameFromAlert(alert string) string {
	s := strings.ToUpper(alert)
	switch {
	case strings.Contains(s, ALGOD):
		return ALGOD
	case strings.Contains(s, MK_PRO):
		return MK_PRO
	case strings.Contains(s, MK_PRO_1):
		return MK_PRO
	case strings.Contains(s, MK_PRO_2):
		return MK_PRO
	case strings.Contains(s, MK_PRO_3):
		return MK_PRO
	case strings.Contains(s, ALGO_PRO_V1):
		return ALGO_PRO_V1
	case strings.Contains(s, ALGO_PRO_V2):
		return ALGO_PRO_V2
	case strings.Contains(s, ALGO_PRO_V3):
		return ALGO_PRO_V3
	case strings.Contains(s, ALGO_SWING):
		return ALGO_SWING
	}
	return NONE
}

func GetIntervalFromAlert(alert string) string {
	s := strings.ToUpper(alert)
	switch {
	case strings.Contains(s, INTERVAL_1M):
		return INTERVAL_1M
	case strings.Contains(s, INTERVAL_5M):
		return INTERVAL_5M
	case strings.Contains(s, INTERVAL_10M):
		return INTERVAL_10M
	case strings.Contains(s, INTERVAL_30M):
		return INTERVAL_30M
	case strings.Contains(s, INTERVAL_1H):
		return INTERVAL_1H
	case strings.Contains(s, INTERVAL_1D):
		return INTERVAL_1D
	}

	return NONE
}

func IsOpenAction(action string) bool {
	return action == BUY || action == BOTTOM
}

func IsCloseAction(action string) bool {
	return action == SELL || action == TOP
}
