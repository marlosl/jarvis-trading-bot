package handlers

import (
	"fmt"
	"jarvis-trading-bot/services/authentication"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type ContextResponse struct {
	Key string `json:"key"`
}

type AuthResponse struct {
	IsAuthorized bool            `json:"isAuthorized"`
	Context      ContextResponse `json:"context"`
}

const AUTHORIZATION = "authorization"

func AuthorizerHandler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request body: %s\n", req.Body)

	auth := req.Headers[AUTHORIZATION]
	auth = strings.Replace(auth, "Secret ", "", 1)
	fmt.Printf("Request Auth: %s\n", auth)

	response := &AuthResponse{
		IsAuthorized: authentication.IsSecretValid(auth),
		Context: ContextResponse{
			Key: "auth-authorizer",
		},
	}
	return ApiResponse(http.StatusOK, response)
}
