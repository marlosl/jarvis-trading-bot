package authentication

import (
	"time"
)

type AuthenticationItem struct {
	Secret    string     `json:"secret" dynamodbav:"secret"`
	CreatedAt *time.Time `json:"createdAt,omitempty" dynamodbav:"createdAt,omitempty"`
}

type AuthenticationDbItem struct {
	PK         string             `json:"pk" dynamodbav:"PK"`
	SK         string             `json:"sk" dynamodbav:"SK"`
	Attributes AuthenticationItem `json:"attributes" dynamodbav:"attributes"`
}
