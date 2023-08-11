package handlers

import (
	"fmt"
	"net/http"
	"os"

	"jarvis-trading-bot/clients/topic"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/signal"
	"jarvis-trading-bot/utils"

	"github.com/aws/aws-lambda-go/events"
)

const ALERT = "alert"

func AlertSignalReceiverHandler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request body: %s\n", req.Body)

	alert := req.PathParameters[ALERT]

	fmt.Printf("Request Alert: %s\n", alert)

	signalItem, err := signal.ConvertAlertToSignal(alert, req.Body)
	if err != nil {
		return ApiResponse(http.StatusBadRequest, fmt.Sprintf("Can't unmarshal body: %v\n", err))
	}

	fmt.Println("Signal: ", utils.SPrintJson(signalItem))

	priceRequestTopic := os.Getenv(consts.SNSPriceRequestTopic)
	snsClient, err := topic.NewSNSClient(&priceRequestTopic)
	if err != nil {
		return ApiResponse(http.StatusInternalServerError, fmt.Sprintf("Can't initialize the topic %s: %v", priceRequestTopic, err))
	}

	err = snsClient.SendMsg(&signalItem)
	if err != nil {
		return ApiResponse(http.StatusInternalServerError, fmt.Sprintf("Can't send the message to the queue: %v", err))
	}
	return ApiResponse(http.StatusOK, signalItem)
}
