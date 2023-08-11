package awsdeploy

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateEventBridgeEvent(ctx *pulumi.Context) error {
	requestPriceScheduleHandlerEventBridgeRule, err := cloudwatch.NewEventRule(ctx, "RequestPriceScheduleHandlerEventBridgeRule", &cloudwatch.EventRuleArgs{
		Description:        pulumi.String("Price request schedule handler event rule"),
		ScheduleExpression: pulumi.String("cron(0,30 * * * ? *)"),
		Name:               pulumi.String("trading-bot-price-request-schedule-rule"),
		IsEnabled:          pulumi.Bool(true),
	})
	if err != nil {
		return err
	}

	_, err = cloudwatch.NewEventTarget(ctx, "RequestPriceScheduleHandlerEventBridgeTarget", &cloudwatch.EventTargetArgs{
		Rule: requestPriceScheduleHandlerEventBridgeRule.Name,
		Arn:  RequestPriceScheduleHandlerLambdaFunction.Arn,
	})
	if err != nil {
		return err
	}

	_, err = lambda.NewPermission(ctx, "RequestPriceScheduleHandlerEventBridgeLambdaPermission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  RequestPriceScheduleHandlerLambdaFunction.Name,
		Principal: pulumi.String("events.amazonaws.com"),
		SourceArn: requestPriceScheduleHandlerEventBridgeRule.Arn,
	})
	if err != nil {
		return err
	}

	return nil
}
