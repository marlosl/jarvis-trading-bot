package tickerconfig

import (
	"fmt"
	"log"
	"os"

	"jarvis-trading-bot/clients/dynamodbclient"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/utils"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type TickerConfigRepository struct {
	dynamodbclient.DBClient
}

var TICKER = "TICKER"

func NewTickerConfigRepository() (*TickerConfigRepository, error) {
	tableName := os.Getenv(consts.ConfigTableName)
	dbClient, err := dynamodbclient.NewDBClient(tableName, nil)
	if err != nil {
		return nil, err
	}

	return &TickerConfigRepository{
		*dbClient,
	}, nil
}

func (db *TickerConfigRepository) SaveItem(item *TickerConfigItem) error {
	svc := dynamodb.New(db.Session)

	if item.CreatedAt == nil || item.CreatedAt.IsZero() {
		now := utils.GetCurrentTime()
		item.CreatedAt = &now
	}

	dbItem := &TickerConfigDbItem{
		PK:         TICKER,
		SK:         fmt.Sprintf("%s#%s", TICKER, item.Ticker),
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

func (db *TickerConfigRepository) GetItem(ticker string) (*TickerConfigItem, error) {
	svc := dynamodb.New(db.Session)

	sk := fmt.Sprintf("%s#%s", TICKER, ticker)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &TICKER,
			},
			"SK": {
				S: &sk,
			},
		},
		TableName: db.TableName,
	}

	result, err := svc.GetItem(input)
	if err != nil {
		fmt.Printf("Got error calling GetItem: %s\n", err)
	}

	dbItem := &TickerConfigDbItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, dbItem)
	if err != nil {
		fmt.Printf("Got error unmarshalling: %s\n", err)
	}

	return &dbItem.Attributes, nil
}

func (db *TickerConfigRepository) DeleteItem(ticker string) error {
	svc := dynamodb.New(db.Session)

	sk := fmt.Sprintf("%s#%s", TICKER, ticker)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &TICKER,
			},
			"SK": {
				S: &sk,
			},
		},
		TableName: db.TableName,
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Printf("Got error calling DeleteItem: %v\n", err)
	}

	fmt.Println("Deleted '" + ticker + "' from table " + *db.TableName)
	return nil
}

func (db *TickerConfigRepository) GetItems() ([]TickerConfigItem, error) {
	svc := dynamodb.New(db.Session)

	keyCond := expression.Key("PK").Equal(expression.Value(TICKER))
	keyCond = keyCond.And(expression.Key("SK").BeginsWith(TICKER))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		Build()

	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 db.TableName,
	}

	result, err := svc.Query(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return nil, err
	}

	items := []TickerConfigItem{}
	for _, i := range result.Items {
		item := TickerConfigDbItem{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
			return nil, err
		}
		items = append(items, item.Attributes)
	}

	return items, nil
}
