package awsdeploy

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const AppPrefix = "trading-bot"

func Deploy(ctx *pulumi.Context) error {
	fmt.Println("Deploying AWS infrastructure...")

	fmt.Println("Creating DynamoDB tables...")
	err := CreateDynamoDbTables(ctx)
	if err != nil {
		fmt.Printf("Can't create DynamoDB tables: %v\n", err)
		return err
	}

	fmt.Println("Creating topics...")
	err = CreateTopics(ctx)
	if err != nil {
		fmt.Printf("Can't create the topics: %v\n", err)
		return err
	}

	fmt.Println("Creating queues...")
	err = CreateQueues(ctx)
	if err != nil {
		fmt.Printf("Can't create the queues: %v\n", err)
		return err
	}

	fmt.Println("Creating log groups...")
	err = CreateLogGroups(ctx)
	if err != nil {
		fmt.Printf("Can't create the log groups: %v\n", err)
		return err
	}

	fmt.Println("Creating roles and policies...")
	err = CreateRolesPolicies(ctx)
	if err != nil {
		fmt.Printf("Can't create the roles and policies: %v\n", err)
		return err
	}

	fmt.Println("Creating lambda functions...")
	err = CreateLambdaFunctions(ctx)
	if err != nil {
		fmt.Printf("Can't create the lambda functions: %v\n", err)
		return err
	}

	fmt.Println("Creating lambda event source mapping...")
	err = CreateLambdaEventSourceMapping(ctx)
	if err != nil {
		fmt.Printf("Can't create the lambda event source mapping: %v\n", err)
		return err
	}

	fmt.Println("Creating EventBridge events...")
	err = CreateEventBridgeEvent(ctx)
	if err != nil {
		fmt.Printf("Can't create EventBridge events: %v\n", err)
		return err
	}

	fmt.Println("Creating the API Gateway...")
	err = CreateApiGateway(ctx)
	if err != nil {
		fmt.Printf("Can't create API Gateway: %v\n", err)
		return err
	}

	return nil
}
