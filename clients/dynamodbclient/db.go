package dynamodbclient

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type DBClient struct {
	TableName *string
	Session   *session.Session
}

func NewDBClient(tableName string, config *aws.Config) (*DBClient, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("Got an error creating a new session:")
		fmt.Println(err)
		return nil, err
	}

	return &DBClient{
		TableName: &tableName,
		Session:   sess,
	}, nil
}
