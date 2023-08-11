package authentication

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

type AuthenticationRepository struct {
	dynamodbclient.DBClient
}

var SECRET = "SECRET"

func NewAuthenticationRepository() (*AuthenticationRepository, error) {
	tableName := os.Getenv(consts.ConfigTableName)
	dbClient, err := dynamodbclient.NewDBClient(tableName, nil)
	if err != nil {
		return nil, err
	}

	return &AuthenticationRepository{
		*dbClient,
	}, nil
}

func (db *AuthenticationRepository) SaveItem(item *AuthenticationItem) error {
	svc := dynamodb.New(db.Session)

	if item.CreatedAt == nil || item.CreatedAt.IsZero() {
		now := utils.GetCurrentTime()
		item.CreatedAt = &now
	}

	dbItem := &AuthenticationDbItem{
		PK:         SECRET,
		SK:         item.Secret,
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

	fmt.Println("Successfully added '" + item.Secret + " to table " + *db.TableName)
	return nil
}

func (db *AuthenticationRepository) GetItem(secret string) (*AuthenticationItem, error) {
	svc := dynamodb.New(db.Session)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &SECRET,
			},
			"SK": {
				S: &secret,
			},
		},
		TableName: db.TableName,
	}

	result, err := svc.GetItem(input)
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	dbItem := &AuthenticationDbItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, dbItem)
	if err != nil {
		log.Fatalf("Got error unmarshalling: %s", err)
	}

	return &dbItem.Attributes, nil
}
