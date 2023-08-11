package handlers

import (
	"fmt"
	"net/http"
	"os"

	"jarvis-trading-bot/clients/topic"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/cache"
	"jarvis-trading-bot/services/signal"

	"github.com/aws/aws-lambda-go/events"
)

func SignalReceiverHandler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request body: %s\n", req.Body)

	signal, err := signal.ConvertTextToSignal(req.Body)
	if err != nil {
		return ApiResponse(http.StatusBadRequest, fmt.Sprintf("Can't unmarshal body: %v\n", err))
	}

	cacheKey := fmt.Sprintf("%s-%s-%s", signal.Ticker, signal.IndicatorName, signal.Interval)

	if cache.Exists(cacheKey, signal.Action) {
		fmt.Printf("Item key %s with value %s already exists in the cache.\n", cacheKey, signal.Action)
		return ApiResponse(http.StatusOK, signal)
	}

	err = cache.Update(cacheKey, signal.Action)
	if err != nil {
		fmt.Printf("Error updating cache: %v\n", err)
	}

	signalTopic := os.Getenv(consts.SNSSignalsTopic)
	snsClient, err := topic.NewSNSClient(&signalTopic)
	if err != nil {
		return ApiResponse(http.StatusInternalServerError, fmt.Sprintf("Can't initialize the topic %s: %v", signalTopic, err))
	}

	err = snsClient.SendMsg(&signal)
	if err != nil {
		return ApiResponse(http.StatusInternalServerError, fmt.Sprintf("Can't send the message to the topic: %v", err))
	}
	return ApiResponse(http.StatusOK, signal)
}
