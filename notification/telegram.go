package notification

import (
	"fmt"
	"net/http"
	"net/url"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"
)

func GetTelegramUrl() string {
	url := utils.GetStringConfig(consts.TelegramUrl)
	botId := utils.GetStringConfig(consts.TelegramBotId)
	apiId := utils.GetStringConfig(consts.TelegramAPIId)

	return url + "/" + botId + ":" + apiId
}

func SendJson(json string) {
	listText := fmt.Sprintf("<code>%s</code>", json)
	SendMessage(listText, true)
}

func SendMessage(message string, isHtml bool) {
	params := url.Values{}
	if isHtml {
		params.Add("parse_mode", "html")
	}
	SendTelegramMessage(message, params)
}

func SendRepliedMessage(message string, reply string) {
	params := url.Values{}
	params.Add("reply_markup", reply)
	SendTelegramMessage(message, params)
}

func SendTelegramMessage(message string, paramValues url.Values) {
	params := url.Values{}
	if (len(paramValues)) > 0 {
		params = paramValues
	}

	chatId := utils.GetStringConfig(consts.TelegramChatId)

	params.Add("chat_id", chatId)
	params.Add("text", message)

	urlMsg := GetTelegramUrl() + "/sendMessage?" + params.Encode()

	log.InfoLogger.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	log.ErrorLogger.Println("err", err)
	log.InfoLogger.Println("response", response)
}

func SendTelegramCallbackQueryResponse(callbackQueryId string) {
	params := url.Values{}

	params.Add("callback_query_id", callbackQueryId)

	urlMsg := GetTelegramUrl() + "/answerCallbackQuery?" + params.Encode()

	log.InfoLogger.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	log.ErrorLogger.Println("err", err)
	log.InfoLogger.Println("response", response)
}
