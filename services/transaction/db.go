package transaction

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"jarvis-trading-bot/clients/dynamodbclient"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/services/types"
	"jarvis-trading-bot/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

const STATUS = "STATUS"
const TICKER_AND_STATUS_INDEX = 1
const STATUS_INDEX = 2

type TransactionRepository struct {
	dynamodbclient.DBClient
}

func NewTransactionRepository(config *aws.Config) (*TransactionRepository, error) {
	tableName := os.Getenv(consts.TransactionsTableName)
	dbClient, err := dynamodbclient.NewDBClient(tableName, config)
	if err != nil {
		return nil, err
	}

	return &TransactionRepository{
		*dbClient,
	}, nil
}

func (db *TransactionRepository) SaveItem(item *types.TransactionItem) error {
	fmt.Printf("Saving item - Ticker: %s, Status: %s", item.Ticker, item.Status)
	svc := dynamodb.New(db.Session)

	if len(item.Uuid) == 0 {
		item.Uuid = uuid.NewString()
	}

	if item.CreatedAt == nil || item.CreatedAt.IsZero() {
		now := utils.GetCurrentTime()
		item.CreatedAt = &now
	}

	dbItem := &TransactionDbItem{
		PK:         item.Ticker,
		SK:         fmt.Sprintf("TRANSACTION#%s#%s", item.CreatedAt.Format("2006-01-02T15:04:05"), item.Uuid),
		GS1PK:      item.Ticker,
		GS1SK:      fmt.Sprintf("%s#%s#%s", item.Status, item.CreatedAt.Format("2006-01-02T15:04:05"), item.Uuid),
		GS2PK:      STATUS,
		GS2SK:      fmt.Sprintf("%s#%s", item.Status, item.Uuid),
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

func (db *TransactionRepository) GetSK(item *types.TransactionItem, status string) string {
	return fmt.Sprintf("%s#%s#%s", status, item.CreatedAt.Format("2006-01-02T15:04:05"), item.Uuid)
}

func (db *TransactionRepository) UpdateItem(item *types.TransactionItem) error {
	fmt.Printf("Updating item - Ticker: %s, Status: %s", item.Ticker, item.Status)
	svc := dynamodb.New(db.Session)

	itemSK := db.GetSK(item, item.Status)
	if item.Status == consts.STATUS_CLOSED {
		itemSK = db.GetSK(item, consts.STATUS_ACTIVE)
	}

	upd := expression.
		Set(expression.Name("PK"), expression.Value(item.Ticker)).
		Set(expression.Name("SK"), expression.Value(db.GetSK(item, item.Status))).
		Set(expression.Name("attributes.signals"), expression.Value(item.Signals)).
		Set(expression.Name("attributes.buyPrice"), expression.Value(item.BuyPrice)).
		Set(expression.Name("attributes.sellPrice"), expression.Value(item.SellPrice)).
		Set(expression.Name("attributes.status"), expression.Value(item.Status))

	expr, err := expression.NewBuilder().WithUpdate(upd).Build()

	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: db.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(item.Ticker),
			},
			"SK": {
				S: aws.String(itemSK),
			},
		},
		ReturnValues:              aws.String("UPDATED_NEW"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	fmt.Println("Successfully updated '" + item.Ticker + "' (" + item.Status + ")")
	return nil
}

func (db *TransactionRepository) CreateExpressionAttributeValues(s interface{}) map[string]*dynamodb.AttributeValue {
	values := make(map[string]*dynamodb.AttributeValue)

	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return values
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("json"), ",")

		fmt.Println("Name:", f.Name)
		for j := 0; j < len(v); j++ {
			fmt.Println(v[j])
		}
	}
	return values
}

func (db *TransactionRepository) GetItems(ticker string, status string) ([]types.TransactionItem, error) {
	svc := dynamodb.New(db.Session)

	keyCond := expression.Key("PK").Equal(expression.Value(ticker))
	keyCond = keyCond.And(expression.Key("SK").BeginsWith(status))

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

	items := []types.TransactionItem{}
	for _, i := range result.Items {
		item := TransactionDbItem{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
			return nil, err
		}
		items = append(items, item.Attributes)
	}

	return items, nil
}

func (db *TransactionRepository) GetItemsByStatus(ticker string, status string) ([]types.TransactionItem, error) {
	return db.GetItemsByIndexedStatus(ticker, status, TICKER_AND_STATUS_INDEX)
}

func (db *TransactionRepository) GetActiveItems() ([]types.TransactionItem, error) {
	return db.GetItemsByIndexedStatus(STATUS, consts.STATUS_ACTIVE, STATUS_INDEX)
}

func (db *TransactionRepository) GetItemsByIndexedStatus(keyValue string, status string, indexType int) ([]types.TransactionItem, error) {

	pkName := ""
	skName := ""
	pkValue := ""
	skValue := ""
	indexName := ""

	switch indexType {
	case TICKER_AND_STATUS_INDEX:
		pkName = "GS1PK"
		skName = "GS1SK"
		pkValue = keyValue
		skValue = status
		indexName = "GS1_INDEX"
	case STATUS_INDEX:
		pkName = "GS2PK"
		skName = "GS2SK"
		pkValue = STATUS
		skValue = status
		indexName = "GS2_INDEX"
	}
	svc := dynamodb.New(db.Session)

	keyCond := expression.Key(pkName).Equal(expression.Value(pkValue))
	keyCond = keyCond.And(expression.Key(skName).BeginsWith(skValue))

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
		IndexName:                 aws.String(indexName),
	}

	result, err := svc.Query(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return nil, err
	}

	items := []types.TransactionItem{}
	for _, i := range result.Items {
		item := TransactionDbItem{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
			return nil, err
		}
		items = append(items, item.Attributes)
	}

	return items, nil
}
