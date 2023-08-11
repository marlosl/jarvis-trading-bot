package main

import (
	"jarvis-trading-bot/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.OperationHandler)
}
