package topic

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSClient struct {
	Topic   *string
	Session *session.Session
}

func NewSNSClient(topic *string) (*SNSClient, error) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Got an error creating a new session:")
		fmt.Println(err)
		return nil, err
	}

	return &SNSClient{
		Topic:   topic,
		Session: sess,
	}, nil
}

func (s *SNSClient) SendMsg(message interface{}) error {
	svc := sns.New(s.Session)

	body, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Got an error marshalling the message:")
		fmt.Println(err)
		return err
	}

	hashBytes := sha1.Sum(body)
	hash := fmt.Sprintf("%x", hashBytes)

	fmt.Println("HASH: ", hash)
	messageSize := len(body)

	_, err = svc.Publish(&sns.PublishInput{
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"MessageSize": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(messageSize)),
			},
		},
		Message:                aws.String(string(body)),
		TopicArn:               s.Topic,
		MessageGroupId:         aws.String("signals"),
		MessageDeduplicationId: aws.String(hash),
	})

	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return err
	}

	fmt.Println("Sent message to the topic")
	return nil
}
