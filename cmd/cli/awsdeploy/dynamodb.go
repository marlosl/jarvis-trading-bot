package awsdeploy

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	SignalsDynamoDbTable      *dynamodb.Table
	TransactionsDynamoDbTable *dynamodb.Table
	ConfigDynamoDbTable       *dynamodb.Table
	CacheDynamoDbTable        *dynamodb.Table
)

func CreateDynamoDbTables(ctx *pulumi.Context) error {
	transactionTable, err := dynamodb.NewTable(ctx, "TransactionsDynamoDbTable", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("SK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("GS1PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("GS1SK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("GS2PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("GS2SK"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		GlobalSecondaryIndexes: dynamodb.TableGlobalSecondaryIndexArray{
			&dynamodb.TableGlobalSecondaryIndexArgs{
				HashKey:        pulumi.String("GS1PK"),
				RangeKey:       pulumi.String("GS1SK"),
				Name:           pulumi.String("GS1_INDEX"),
				ProjectionType: pulumi.String("ALL"),
			},
			&dynamodb.TableGlobalSecondaryIndexArgs{
				HashKey:        pulumi.String("GS2PK"),
				RangeKey:       pulumi.String("GS2SK"),
				Name:           pulumi.String("GS2_INDEX"),
				ProjectionType: pulumi.String("ALL"),
			},
		},
		HashKey:  pulumi.String("PK"),
		RangeKey: pulumi.String("SK"),
		Name:     pulumi.String("trading-bot-transactions"),
	})
	if err != nil {
		return err
	}
	TransactionsDynamoDbTable = transactionTable
	ctx.Export("transactionTable", transactionTable.Arn)

	signalTable, err := dynamodb.NewTable(ctx, "SignalsDynamoDbTable", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("SK"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("PK"),
		RangeKey:    pulumi.String("SK"),
		Name:        pulumi.String("trading-bot-signals"),
	})
	if err != nil {
		return err
	}
	SignalsDynamoDbTable = signalTable
	ctx.Export("signalTable", signalTable.Arn)

	configTable, err := dynamodb.NewTable(ctx, "ConfigDynamoDbTable", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("SK"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("PK"),
		RangeKey:    pulumi.String("SK"),
		Name:        pulumi.String("trading-bot-config"),
	})
	if err != nil {
		return err
	}
	ConfigDynamoDbTable = configTable
	ctx.Export("configTable", configTable.Arn)

    cacheTable, err := dynamodb.NewTable(ctx, "CacheDynamoDbTable", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("SK"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("PK"),
		RangeKey:    pulumi.String("SK"),
		Name:        pulumi.String("trading-bot-cache"),
	})
	if err != nil {
		return err
	}
	CacheDynamoDbTable = cacheTable
	ctx.Export("cacheTable", cacheTable.Arn)

	return nil
}
