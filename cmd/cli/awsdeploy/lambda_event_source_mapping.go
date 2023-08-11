package awsdeploy

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateLambdaEventSourceMapping(ctx *pulumi.Context) error {
	processSignalHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping, err := lambda.NewEventSourceMapping(ctx, "ProcessSignalHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping", &lambda.EventSourceMappingArgs{
		EventSourceArn:   SQSSignalsQueue.Arn,
		FunctionName:     ProcessSignalHandlerLambdaFunction.Arn,
		BatchSize:        pulumi.Int(1),
		Enabled:          pulumi.Bool(true),
	},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution}),
	)
	if err != nil {
		return err
	}
	ctx.Export("ProcessSignalHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping", processSignalHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping.ID())

	signalAnalyserHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping, err := lambda.NewEventSourceMapping(ctx, "SignalAnalyserHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping", &lambda.EventSourceMappingArgs{
		EventSourceArn:   SQSSignalsAnalyserQueue.Arn,
		FunctionName:     SignalAnalyserHandlerLambdaFunction.Arn,
		BatchSize:        pulumi.Int(1),
		Enabled:          pulumi.Bool(true),
	},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution}),
	)
	if err != nil {
		return err
	}
	ctx.Export("SignalAnalyserHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping", signalAnalyserHandlerLambdaFunctionSQSSignalsQueueEventSourceMapping.ID())

	priceRequestHandlerLambdaFunctionSQSPriceRequestQueueEventSourceMapping, err := lambda.NewEventSourceMapping(ctx, "PriceRequestHandlerLambdaFunctionSQSPriceRequestQueueEventSourceMapping", &lambda.EventSourceMappingArgs{
		EventSourceArn:   SQSPriceRequestQueue.Arn,
		FunctionName:     PriceRequestHandlerLambdaFunction.Arn,
		BatchSize:        pulumi.Int(1),
		Enabled:          pulumi.Bool(true),
	},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution}),
	)
	if err != nil {
		return err
	}
	ctx.Export("PriceRequestHandlerLambdaFunctionSQSPriceRequestQueueEventSourceMapping", priceRequestHandlerLambdaFunctionSQSPriceRequestQueueEventSourceMapping.ID())
	
	operationHandlerLambdaFunctionSQSOperationEventSourceMapping, err := lambda.NewEventSourceMapping(ctx, "OperationHandlerLambdaFunctionSQSOperationEventSourceMapping", &lambda.EventSourceMappingArgs{
		EventSourceArn:   SQSOperationQueue.Arn,
		FunctionName:     OperationHandlerLambdaFunction.Arn,
		BatchSize:        pulumi.Int(1),
		Enabled:          pulumi.Bool(true),
	},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution}),
	)
	if err != nil {
		return err
	}
	ctx.Export("OperationHandlerLambdaFunctionSQSOperationEventSourceMapping", operationHandlerLambdaFunctionSQSOperationEventSourceMapping.ID())

	return nil
}
