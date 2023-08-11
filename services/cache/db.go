package cache

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

var (
	CACHE = "CACHE"
)

type CacheRepository struct {
	dynamodbclient.DBClient
}

func NewCacheRepository() (*CacheRepository, error) {
	tableName := os.Getenv(consts.CacheTableName)
	dbClient, err := dynamodbclient.NewDBClient(tableName, nil)
	if err != nil {
		return nil, err
	}

	return &CacheRepository{
		*dbClient,
	}, nil
}

func (db *CacheRepository) SaveItem(key, value string) error {
	svc := dynamodb.New(db.Session)

	dbItem := &CacheItem{
		PK:    CACHE,
		SK:    key,
		Value: value,
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

	fmt.Println("Successfully added '" + key + " to table " + *db.TableName)
	return nil
}

func (db *CacheRepository) GetItem(key string) (*CacheItem, error) {
	svc := dynamodb.New(db.Session)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &CACHE,
			},
			"SK": {
				S: &key,
			},
		},
		TableName: db.TableName,
	}

	result, err := svc.GetItem(input)
	if err != nil {
		fmt.Printf("Got error calling GetItem: %s\n", err)
	}

	dbItem := &CacheItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, dbItem)
	if err != nil {
		fmt.Printf("Got error unmarshalling: %s\n", err)
	}

	return dbItem, nil
}

func (db *CacheRepository) DeleteItem(key string) error {
	svc := dynamodb.New(db.Session)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &CACHE,
			},
			"SK": {
				S: &key,
			},
		},
		TableName: db.TableName,
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Printf("Got error calling DeleteItem: %v\n", err)
	}

	fmt.Println("Deleted '" + key + "' from table " + *db.TableName)
	return nil
}
