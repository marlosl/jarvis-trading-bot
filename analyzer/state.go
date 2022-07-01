package analyzer

import (
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/utils"
)

func GetCandles(symbol string) []structs.Candlestick {
	var candles []structs.Candlestick = make([]structs.Candlestick, 0)
	utils.DB.Where("symbol = ?", symbol).Order("event_time").Find(&candles)
	return candles
}

func GetPendingOrders(intanceId string, symbol string) []structs.Operation {
	var operations []structs.Operation = make([]structs.Operation, 0)
	utils.DB.
		Where("instance_id = ? AND symbol = ? AND finished = false", intanceId, symbol).
		Order("transact_time").
		Find(&operations)
	return operations
}

func GetCryptoBotParameters() []structs.BotParameters {
	var params []structs.BotParameters = make([]structs.BotParameters, 0)
	utils.DB.Where("closed is null or closed < to_timestamp(0)").Find(&params)
	return params
}

func GetCryptoBotParametersByUser(userId string) []structs.BotParameters {
	var params []structs.BotParameters = make([]structs.BotParameters, 0)
	utils.DB.Where("user_id = ? and closed is null or closed < to_timestamp(0)", utils.ConvertStringToInt(userId)).Find(&params)
	return params
}

func GetCryptoBotParametersById(id string) *structs.BotParameters {
	params := new(structs.BotParameters)
	utils.DB.Where("id = ?", utils.ConvertStringToInt(id)).Find(&params)
	return params
}

func GetTradingStatus(botParametersId uint, userId uint, simulation bool) *[]structs.TradingStatus {
	var ts []structs.TradingStatus = make([]structs.TradingStatus, 0)
	utils.DB.Where("bot_param_id =? AND user_id = ? AND simulation = ? AND (deleted_at is null or deleted_at < to_timestamp(0))", botParametersId, userId, simulation).Find(&ts)
	if len(ts) > 0 {
		return &ts
	}
	return nil
}

func GetCandlesticks(symbol string, startDate string, endDate string) *[]structs.Candlestick {
	var ts []structs.Candlestick = make([]structs.Candlestick, 0)
	utils.DB.Where("symbol =? AND (event_time >= ? AND event_time <= ?)", symbol, startDate, endDate).Order("event_time").Find(&ts)
	if len(ts) > 0 {
		return &ts
	}
	return nil
}
