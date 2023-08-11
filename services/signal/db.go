package signal

import (
	"fmt"
	"log"
	"os"

	"jarvis-trading-bot/clients/dynamodbclient"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/utils"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type SignalRepository struct {
	dynamodbclient.DBClient
}

func NewSignalRepository() (*SignalRepository, error) {
	tableName := os.Getenv(consts.SignalsTableName)
	dbClient, err := dynamodbclient.NewDBClient(tableName, nil)
	if err != nil {
		return nil, err
	}

	return &SignalRepository{
		*dbClient,
	}, nil
}

func (db *SignalRepository) SaveItem(item *SignalItem) error {
	svc := dynamodb.New(db.Session)

	if item.CreatedAt == nil || item.CreatedAt.IsZero() {
		now := utils.GetCurrentTime()
		item.CreatedAt = &now
	}

	dbItem := &SignalDbItem{
		PK:         item.Ticker,
		SK:         fmt.Sprintf("%s#%s#%s", item.Action, item.IndicatorName, uuid.NewString()),
		Attributes: *item,
	}

	fmt.Printf("Saving item: %s", utils.SPrintJson(dbItem))

	av, err := dynamodbattribute.MarshalMap(dbItem)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: db.TableName,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + item.Ticker + " to table " + *db.TableName)
	return nil
}
