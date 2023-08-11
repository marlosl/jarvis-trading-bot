package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"jarvis-trading-bot/services/exchangeconfig"

	"github.com/aws/aws-lambda-go/events"
)

func ExchangeConfigHandler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request body: %s\n", req.Body)

	var config exchangeconfig.ExchangeConfigItem

	err := json.Unmarshal([]byte(req.Body), &config)
	if err != nil {
		return ApiResponse(http.StatusBadRequest, fmt.Sprintf("Can't unmarshal body: %v\n", err))
	}

	err = exchangeconfig.SaveExchangeConfig(&config)
	if err != nil {
		return ApiResponse(http.StatusInternalServerError, fmt.Sprintf("Can't save the exchange config: %v", err))
	}
	return ApiResponse(http.StatusOK, config)
}
