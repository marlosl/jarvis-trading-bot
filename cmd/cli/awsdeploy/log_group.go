package awsdeploy

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	AuthorizerHandlerLogGroup           *cloudwatch.LogGroup
	ReceiveSignalHandlerLogGroup        *cloudwatch.LogGroup
	ReceiveAlertSignalHandlerLogGroup   *cloudwatch.LogGroup
	ProcessSignalHandlerLogGroup        *cloudwatch.LogGroup
	SignalAnalyserHandlerLogGroup       *cloudwatch.LogGroup
	PriceRequestHandlerLogGroup         *cloudwatch.LogGroup
	ExchangeConfigHandlerLogGroup       *cloudwatch.LogGroup
	PriceRequestScheduleHandlerLogGroup *cloudwatch.LogGroup
	ApiGatewayLogGroup                  *cloudwatch.LogGroup
	OperationHandlerLogGroup            *cloudwatch.LogGroup
)

func CreateLogGroups(ctx *pulumi.Context) error {
	authorizerHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "AuthorizerHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-authorizer-handler"),
	})

	if err != nil {
		return err
	}
	AuthorizerHandlerLogGroup = authorizerHandlerLogGroup

	receiveSignalHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "ReceiveSignalHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-receive-signal-handler"),
	})

	if err != nil {
		return err
	}
	ReceiveSignalHandlerLogGroup = receiveSignalHandlerLogGroup

	receiveAlertSignalHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "ReceiveAlertSignalHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-receive-alert-signal-handler"),
	})

	if err != nil {
		return err
	}
	ReceiveAlertSignalHandlerLogGroup = receiveAlertSignalHandlerLogGroup

	processSignalHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "ProcessSignalHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-process-signal-handler"),
	})

	if err != nil {
		return err
	}
	ProcessSignalHandlerLogGroup = processSignalHandlerLogGroup

	signalAnalyserHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "SignalAnalyserHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-signal-analyser-handler"),
	})

	if err != nil {
		return err
	}
	SignalAnalyserHandlerLogGroup = signalAnalyserHandlerLogGroup

	priceRequestHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "PriceRequestHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-price-request-handler"),
	})

	if err != nil {
		return err
	}
	PriceRequestHandlerLogGroup = priceRequestHandlerLogGroup

	exchangeConfigHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "ExchangeConfigHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-exchange-config-handler"),
	})

	if err != nil {
		return err
	}
	ExchangeConfigHandlerLogGroup = exchangeConfigHandlerLogGroup

	priceRequestScheduleHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "PriceRequestScheduleHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-price-request-schedule-handler"),
	})

	if err != nil {
		return err
	}
	PriceRequestScheduleHandlerLogGroup = priceRequestScheduleHandlerLogGroup

	apiGatewayLogGroup, err := cloudwatch.NewLogGroup(ctx, "ApiGatewayLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/api-gateway/trading-bot"),
	})

	if err != nil {
		return err
	}
	ApiGatewayLogGroup = apiGatewayLogGroup

	operationHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "Operatio nHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/trading-bot-operation-handler"),
	})

	if err != nil {
		return err
	}
	OperationHandlerLogGroup = operationHandlerLogGroup

	return nil
}
