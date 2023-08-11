package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func ProcessRecordsHandler(_ context.Context, kinesesEvent events.KinesisFirehoseEvent) error {
	for _, record := range kinesesEvent.Records {
		fmt.Println(record)
	}
	return nil
}
