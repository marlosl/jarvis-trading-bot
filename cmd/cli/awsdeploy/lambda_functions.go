package awsdeploy

import (
	"path/filepath"

	"jarvis-trading-bot/consts"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	ProcessSignalHandlerLambdaFunction        *lambda.Function
	SignalAnalyserHandlerLambdaFunction       *lambda.Function
	PriceRequestHandlerLambdaFunction         *lambda.Function
	RequestPriceScheduleHandlerLambdaFunction *lambda.Function
	AuthorizerHandlerLambdaFunction           *lambda.Function
	ExchangeConfigHandlerLambdaFunction       *lambda.Function
	ReceiveSignalHandlerLambdaFunction        *lambda.Function
	ReceiveAlertSignalHandlerLambdaFunction   *lambda.Function
	OperationHandlerLambdaFunction            *lambda.Function
)

func CreateLambdaFunctions(ctx *pulumi.Context) error {

	outputDir := os.Getenv(consts.ProjectOutputDir)

	authorizerFile := filepath.Join(outputDir, "authorizer/authorizer-handler.zip")
	authorizerHandlerLambdaFunction, err := lambda.NewFunction(ctx, "AuthorizerHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-authorizer-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(authorizerFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION": pulumi.String(os.Getenv(consts.AwsRegion)),
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, AuthorizerHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	AuthorizerHandlerLambdaFunction = authorizerHandlerLambdaFunction
	ctx.Export("authorizerHandlerLambdaFunction", authorizerHandlerLambdaFunction.Arn)

	signalReceiverFile := filepath.Join(outputDir, "signalreceiver/receive-signal-handler.zip")
	receiveSignalHandlerLambdaFunction, err := lambda.NewFunction(ctx, "ReceiveSignalHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-receive-signal-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(signalReceiverFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":                  pulumi.String(os.Getenv(consts.AwsRegion)),
				"SIGNALS_TABLE_NAME":      SignalsDynamoDbTable.Name,
				"TRANSACTIONS_TABLE_NAME": TransactionsDynamoDbTable.Name,
				"CACHE_TABLE_NAME":      CacheDynamoDbTable.Name,
				"SQS_SIGNALS_QUEUE_URL":   SQSSignalsQueue.Url,
				"SNS_SIGNALS_TOPIC":       SNSSignalsTopic,
				"SNS_PRICE_REQUEST_TOPIC": SNSPriceRequestTopic,
				"SNS_OPERATION_TOPIC":     SNSOperationTopic,
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, ReceiveSignalHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	ReceiveSignalHandlerLambdaFunction = receiveSignalHandlerLambdaFunction
	ctx.Export("ReceiveSignalHandlerLambdaFunction", receiveSignalHandlerLambdaFunction.Arn)

	priceRequestSchedulerFile := filepath.Join(outputDir, "pricerequestscheduler/price-request-schedule-handler.zip")
	requestPriceScheduleHandlerLambdaFunction, err := lambda.NewFunction(ctx, "RequestPriceScheduleHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-price-request-schedule-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(priceRequestSchedulerFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":                  pulumi.String(os.Getenv(consts.AwsRegion)),
				"SIGNALS_TABLE_NAME":      SignalsDynamoDbTable.Name,
				"TRANSACTIONS_TABLE_NAME": TransactionsDynamoDbTable.Name,
				"CACHE_TABLE_NAME":      CacheDynamoDbTable.Name,
				"SQS_SIGNALS_QUEUE_URL":   SQSSignalsQueue.Url,
				"SNS_SIGNALS_TOPIC":       SNSSignalsTopic,
				"SNS_PRICE_REQUEST_TOPIC": SNSPriceRequestTopic,
				"SNS_OPERATION_TOPIC":     SNSOperationTopic,
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, PriceRequestScheduleHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	RequestPriceScheduleHandlerLambdaFunction = requestPriceScheduleHandlerLambdaFunction
	ctx.Export("RequestPriceScheduleHandlerLambdaFunction", requestPriceScheduleHandlerLambdaFunction.Arn)

	signalAnalyserFile := filepath.Join(outputDir, "signalanalyser/signal-analyser-handler.zip")
	signalAnalyserHandlerLambdaFunction, err := lambda.NewFunction(ctx, "SignalAnalyserHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-signal-analyser-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(signalAnalyserFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":                  pulumi.String(os.Getenv(consts.AwsRegion)),
				"SIGNALS_TABLE_NAME":      SignalsDynamoDbTable.Name,
				"TRANSACTIONS_TABLE_NAME": TransactionsDynamoDbTable.Name,
				"CACHE_TABLE_NAME":      CacheDynamoDbTable.Name,
				"CONFIG_TABLE_NAME":       ConfigDynamoDbTable.Name,
				"SQS_SIGNALS_QUEUE_URL":   SQSSignalsQueue.Url,
				"SNS_SIGNALS_TOPIC":       SNSSignalsTopic,
				"SNS_PRICE_REQUEST_TOPIC": SNSPriceRequestTopic,
				"SNS_OPERATION_TOPIC":     SNSOperationTopic,
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, SignalAnalyserHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	SignalAnalyserHandlerLambdaFunction = signalAnalyserHandlerLambdaFunction
	ctx.Export("SignalAnalyserHandlerLambdaFunction", signalAnalyserHandlerLambdaFunction.Arn)

	alertSignalReceiverFile := filepath.Join(outputDir, "alertsignalreceiver/receive-alert-signal-handler.zip")
	receiveAlertSignalHandlerLambdaFunction, err := lambda.NewFunction(ctx, "ReceiveAlertSignalHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-receive-alert-signal-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(alertSignalReceiverFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":                  pulumi.String(os.Getenv(consts.AwsRegion)),
				"SIGNALS_TABLE_NAME":      SignalsDynamoDbTable.Name,
				"TRANSACTIONS_TABLE_NAME": TransactionsDynamoDbTable.Name,
				"CACHE_TABLE_NAME":      CacheDynamoDbTable.Name,
				"CONFIG_TABLE_NAME":       ConfigDynamoDbTable.Name,
				"SQS_SIGNALS_QUEUE_URL":   SQSSignalsQueue.Url,
				"SNS_SIGNALS_TOPIC":       SNSSignalsTopic,
				"SNS_PRICE_REQUEST_TOPIC": SNSPriceRequestTopic,
				"SNS_OPERATION_TOPIC":     SNSOperationTopic,
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, ReceiveAlertSignalHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	ReceiveAlertSignalHandlerLambdaFunction = receiveAlertSignalHandlerLambdaFunction
	ctx.Export("ReceiveAlertSignalHandlerLambdaFunction", receiveAlertSignalHandlerLambdaFunction.Arn)

	signalProcessorFile := filepath.Join(outputDir, "signalprocessor/process-signal-handler.zip")
	processSignalHandlerLambdaFunction, err := lambda.NewFunction(ctx, "ProcessSignalHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-process-signal-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(signalProcessorFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":                  pulumi.String(os.Getenv(consts.AwsRegion)),
				"SIGNALS_TABLE_NAME":      SignalsDynamoDbTable.Name,
				"TRANSACTIONS_TABLE_NAME": TransactionsDynamoDbTable.Name,
				"CACHE_TABLE_NAME":      CacheDynamoDbTable.Name,
				"CONFIG_TABLE_NAME":       ConfigDynamoDbTable.Name,
				"SQS_SIGNALS_QUEUE_URL":   SQSSignalsQueue.Url,
				"SNS_SIGNALS_TOPIC":       SNSSignalsTopic,
				"SNS_PRICE_REQUEST_TOPIC": SNSPriceRequestTopic,
				"SNS_OPERATION_TOPIC":     SNSOperationTopic,
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, ProcessSignalHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	ProcessSignalHandlerLambdaFunction = processSignalHandlerLambdaFunction
	ctx.Export("ProcessSignalHandlerLambdaFunction", processSignalHandlerLambdaFunction.Arn)

	priceRequestFile := filepath.Join(outputDir, "pricerequest/price-request-handler.zip")
	priceRequestHandlerLambdaFunction, err := lambda.NewFunction(ctx, "PriceRequestHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-price-request-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(priceRequestFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":                  pulumi.String(os.Getenv(consts.AwsRegion)),
				"SIGNALS_TABLE_NAME":      SignalsDynamoDbTable.Name,
				"TRANSACTIONS_TABLE_NAME": TransactionsDynamoDbTable.Name,
				"CACHE_TABLE_NAME":      CacheDynamoDbTable.Name,
				"CONFIG_TABLE_NAME":       ConfigDynamoDbTable.Name,
				"SQS_SIGNALS_QUEUE_URL":   SQSSignalsQueue.Url,
				"SNS_SIGNALS_TOPIC":       SNSSignalsTopic,
				"SNS_PRICE_REQUEST_TOPIC": SNSPriceRequestTopic,
				"BINANCE_URL":             pulumi.String("https://data.binance.com"),
				"SNS_OPERATION_TOPIC":     SNSOperationTopic,
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, PriceRequestHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	PriceRequestHandlerLambdaFunction = priceRequestHandlerLambdaFunction
	ctx.Export("PriceRequestHandlerLambdaFunction", priceRequestHandlerLambdaFunction.Arn)

	exchangeConfigFile := filepath.Join(outputDir, "exchangeconfig/exchange-config-handler.zip")
	exchangeConfigHandlerLambdaFunction, err := lambda.NewFunction(ctx, "ExchangeConfigHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-exchange-config-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(exchangeConfigFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":                  pulumi.String(os.Getenv(consts.AwsRegion)),
				"SIGNALS_TABLE_NAME":      SignalsDynamoDbTable.Name,
				"TRANSACTIONS_TABLE_NAME": TransactionsDynamoDbTable.Name,
				"CACHE_TABLE_NAME":      CacheDynamoDbTable.Name,
				"CONFIG_TABLE_NAME":       ConfigDynamoDbTable.Name,
				"SQS_SIGNALS_QUEUE_URL":   SQSSignalsQueue.Url,
				"SNS_SIGNALS_TOPIC":       SNSSignalsTopic,
				"SNS_PRICE_REQUEST_TOPIC": SNSPriceRequestTopic,
				"SNS_OPERATION_TOPIC":     SNSOperationTopic,
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, ExchangeConfigHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	ExchangeConfigHandlerLambdaFunction = exchangeConfigHandlerLambdaFunction
	ctx.Export("ExchangeConfigHandlerLambdaFunction", exchangeConfigHandlerLambdaFunction.Arn)

	operationFile := filepath.Join(outputDir, "operation/operation-handler.zip")
	operationHandlerLambdaFunction, err := lambda.NewFunction(ctx, "OperationHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("trading-bot-operation-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(operationFile),
		Timeout:    pulumi.Int(10),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION": pulumi.String(os.Getenv(consts.AwsRegion)),
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, OperationHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	OperationHandlerLambdaFunction = operationHandlerLambdaFunction
	ctx.Export("OperationHandlerLambdaFunction", operationHandlerLambdaFunction.Arn)

	return nil
}
