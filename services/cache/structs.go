package cache

type CacheItem struct {
	PK         string     `json:"pk" dynamodbav:"PK"`
	SK         string     `json:"sk" dynamodbav:"SK"`
	Value      string     `json:"value" dynamodbav:"Value"`
}
