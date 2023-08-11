package handlers

import (
	"context"
	"fmt"

	"jarvis-trading-bot/services/pricerequestscheduler"
)

func PriceRequestScheduleHandler(_ context.Context, _ interface{}) error {
	err := pricerequestscheduler.CreatePriceRequests()
	if err != nil {
		fmt.Printf("Error creating prices: %v\n", err)
		return err
	}
	return nil
}
