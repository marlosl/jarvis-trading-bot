package exchangeconfig

import (
	"fmt"
	"log"
	"os"

	"jarvis-trading-bot/clients/dynamodbclient"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/utils"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ExchangeConfigRepository struct {
	dynamodbclient.DBClient
}

func NewExchangeConfigRepository() (*ExchangeConfigRepository, error) {
	tableName := os.Getenv(consts.ConfigTableName)
	dbClient, err := dynamodbclient.NewDBClient(tableName, nil)
	if err != nil {
		return nil, err
	}

	return &ExchangeConfigRepository{
		*dbClient,
	}, nil
}

func (db *ExchangeConfigRepository) SaveItem(item *ExchangeConfigItem) error {
	svc := dynamodb.New(db.Session)

	if item.CreatedAt == nil || item.CreatedAt.IsZero() {
		now := utils.GetCurrentTime()
		item.CreatedAt = &now
	}

	dbItem := &ExchangeConfigDbItem{
		PK:         item.Ticker,
		SK:         item.Exchange,
		Attributes: *item,
	}

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

func (db *ExchangeConfigRepository) GetItem(ticker string, exchange string) (*ExchangeConfigItem, error) {
	svc := dynamodb.New(db.Session)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &ticker,
			},
			"SK": {
				S: &exchange,
			},
		},
		TableName: db.TableName,
	}

	result, err := svc.GetItem(input)
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	dbItem := &ExchangeConfigDbItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, dbItem)
	if err != nil {
		log.Fatalf("Got error unmarshalling: %s", err)
	}

	return &dbItem.Attributes, nil
}
