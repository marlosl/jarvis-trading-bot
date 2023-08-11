package awsdeploy

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	SQSSignalsQueue         *sqs.Queue
	SQSSignalsAnalyserQueue *sqs.Queue
	SQSPriceRequestQueue    *sqs.Queue
	SQSOperationQueue       *sqs.Queue
)

func CreateQueues(ctx *pulumi.Context) error {
	err := CreateSignalsQueue(ctx)
	if err != nil {
		return err
	}

	err = CreateSignalsAnalyserQueue(ctx)
	if err != nil {
		return err
	}

	err = CreatePriceRequestQueue(ctx)
	if err != nil {
		return err
	}

	err = CreateOperationQueue(ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateSignalsQueue(ctx *pulumi.Context) error {
	sqsSignalsDeadLetterQueue, err := sqs.NewQueue(ctx, "SQSSignalsDeadLetterQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-signals-dlq.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsSignalsDeadLetterQueue", sqsSignalsDeadLetterQueue.Arn)

	sqsSignalsQueue, err := sqs.NewQueue(ctx, "SQSSignalsQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-signals.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
		RedrivePolicy: pulumi.Sprintf(`{
			"deadLetterTargetArn": "%s",
			"maxReceiveCount": 5
		}`, sqsSignalsDeadLetterQueue.Arn),
	})

	if err != nil {
		return err
	}

	SQSSignalsQueue = sqsSignalsQueue
	ctx.Export("sqsSignalsDeadLetterQueue", sqsSignalsQueue.Arn)

	sqsSignalsQueueSubscription, err := sns.NewTopicSubscription(ctx, "SQSSignalsQueueSubscription", &sns.TopicSubscriptionArgs{
		Topic:    SNSSignalsTopic,
		Protocol: pulumi.String("sqs"),
		Endpoint: sqsSignalsQueue.Arn,
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsSignalsQueueSubscription", sqsSignalsQueueSubscription.Arn)

	sqsSignalsQueuePolicy, err := sqs.NewQueuePolicy(ctx, "SQSSignalsQueuePolicy", &sqs.QueuePolicyArgs{
		QueueUrl: sqsSignalsQueue.ID(),
		Policy: pulumi.Sprintf(`{
			"Version": "2012-10-17",
			"Id": "sqs-signals-queue-policy",
			"Statement": [
				{
					"Sid": "sqs-signals-queue-policy",
					"Effect": "Allow",
					"Principal": "*",
					"Action": "sqs:SendMessage",
					"Resource": "%v",
					"Condition": {
						"ArnEquals": {
							"aws:SourceArn": "%v"
						}
					}
				}
			]
		}`, sqsSignalsQueue.Arn, SNSSignalsTopic),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsSignalsQueuePolicy", sqsSignalsQueuePolicy.ID())

	return nil
}

func CreateSignalsAnalyserQueue(ctx *pulumi.Context) error {
	sqsSignalsAnalyserDeadLetterQueue, err := sqs.NewQueue(ctx, "SQSSignalsAnalyserDeadLetterQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-signals-analyser-dlq.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsSignalsAnalyserDeadLetterQueue", sqsSignalsAnalyserDeadLetterQueue.Arn)

	sqsSignalsAnalyserQueue, err := sqs.NewQueue(ctx, "SQSSignalsAnalyserQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-signals-analyser.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
		RedrivePolicy: pulumi.Sprintf(`{
			"deadLetterTargetArn": "%s",
			"maxReceiveCount": 5
		}`, sqsSignalsAnalyserDeadLetterQueue.Arn),
	})

	if err != nil {
		return err
	}

	SQSSignalsAnalyserQueue = sqsSignalsAnalyserQueue
	ctx.Export("sqsSignalsAnalyserQueue", sqsSignalsAnalyserQueue.Arn)

	sqsSignalsQueueSubscription, err := sns.NewTopicSubscription(ctx, "SQSSignalsAnalyserQueueSubscription", &sns.TopicSubscriptionArgs{
		Topic:    SNSSignalsTopic,
		Protocol: pulumi.String("sqs"),
		Endpoint: sqsSignalsAnalyserQueue.Arn,
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsSignalsQueueSubscription", sqsSignalsQueueSubscription.Arn)

	sqsSignalsAnalyserQueuePolicy, err := sqs.NewQueuePolicy(ctx, "SQSSignalsAnalyserQueuePolicy", &sqs.QueuePolicyArgs{
		QueueUrl: sqsSignalsAnalyserQueue.ID(),
		Policy: pulumi.Sprintf(`{
			"Version": "2012-10-17",
			"Id": "sqs-signals-analyser-policy",
			"Statement": [
				{
					"Sid": "sqs-signals-analyser-policy",
					"Effect": "Allow",
					"Principal": "*",
					"Action": "sqs:SendMessage",
					"Resource": "%v",
					"Condition": {
						"ArnEquals": {
							"aws:SourceArn": "%v"
						}
					}
				}
			]
		}`, sqsSignalsAnalyserQueue.Arn, SNSSignalsTopic),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsSignalsAnalyserQueuePolicy", sqsSignalsAnalyserQueuePolicy.ID())

	return nil
}

func CreatePriceRequestQueue(ctx *pulumi.Context) error {
	sqsPriceRequestDeadLetterQueue, err := sqs.NewQueue(ctx, "SQSPriceRequestDeadLetterQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-price-request-dlq.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsPriceRequestDeadLetterQueue", sqsPriceRequestDeadLetterQueue.Arn)

	sqsPriceRequestQueue, err := sqs.NewQueue(ctx, "SQSPriceRequestQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-price-request.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
		RedrivePolicy: pulumi.Sprintf(`{
			"deadLetterTargetArn": "%s",
			"maxReceiveCount": 5
		}`, sqsPriceRequestDeadLetterQueue.Arn),
	})

	if err != nil {
		return err
	}

	SQSPriceRequestQueue = sqsPriceRequestQueue
	ctx.Export("sqsPriceRequestQueue", sqsPriceRequestQueue.Arn)

	sqsPriceRequestQueueSubscription, err := sns.NewTopicSubscription(ctx, "SQSPriceRequestQueueSubscription", &sns.TopicSubscriptionArgs{
		Topic:    SNSSignalsTopic,
		Protocol: pulumi.String("sqs"),
		Endpoint: sqsPriceRequestQueue.Arn,
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsPriceRequestQueueSubscription", sqsPriceRequestQueueSubscription.Arn)

	sqsPriceRequestQueuePolicy, err := sqs.NewQueuePolicy(ctx, "SQSPriceRequestQueuePolicy", &sqs.QueuePolicyArgs{
		QueueUrl: sqsPriceRequestQueue.ID(),
		Policy: pulumi.Sprintf(`{
			"Version": "2012-10-17",
			"Id": "sqs-price-request-policy",
			"Statement": [
				{
					"Sid": "sqs-price-request-policy",
					"Effect": "Allow",
					"Principal": "*",
					"Action": "sqs:SendMessage",
					"Resource": "%v",
					"Condition": {
						"ArnEquals": {
							"aws:SourceArn": "%v"
						}
					}
				}
			]
		}`, sqsPriceRequestQueue.Arn, SNSPriceRequestTopic),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsPriceRequestQueuePolicy", sqsPriceRequestQueuePolicy.ID())

	return nil
}

func CreateOperationQueue(ctx *pulumi.Context) error {
	sqsOperationDeadLetterQueue, err := sqs.NewQueue(ctx, "SQSOperationDeadLetterQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-operation-dlq.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsOperationDeadLetterQueue", sqsOperationDeadLetterQueue.Arn)

	sqsOperationQueue, err := sqs.NewQueue(ctx, "SQSOperationQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("trading-bot-operation.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
		RedrivePolicy: pulumi.Sprintf(`{
			"deadLetterTargetArn": "%s",
			"maxReceiveCount": 5
		}`, sqsOperationDeadLetterQueue.Arn),
	})

	if err != nil {
		return err
	}

	SQSOperationQueue = sqsOperationQueue
	ctx.Export("sqsOperationQueue", sqsOperationQueue.Arn)

	sqsOperationQueueSubscription, err := sns.NewTopicSubscription(ctx, "SQSOperationQueueSubscription", &sns.TopicSubscriptionArgs{
		Topic:    SNSSignalsTopic,
		Protocol: pulumi.String("sqs"),
		Endpoint: sqsOperationQueue.Arn,
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsOperationQueueSubscription", sqsOperationQueueSubscription.Arn)

	sqsOperationQueuePolicy, err := sqs.NewQueuePolicy(ctx, "SQSOperationQueuePolicy", &sqs.QueuePolicyArgs{
		QueueUrl: sqsOperationQueue.ID(),
		Policy: pulumi.Sprintf(`{
			"Version": "2012-10-17",
			"Id": "sqs-operation-policy",
			"Statement": [
				{
					"Sid": "sqs-operation-policy",
					"Effect": "Allow",
					"Principal": "*",
					"Action": "sqs:SendMessage",
					"Resource": "%v",
					"Condition": {
						"ArnEquals": {
							"aws:SourceArn": "%v"
						}
					}
				}
			]
		}`, sqsOperationQueue.Arn, SNSOperationTopic),
	})

	if err != nil {
		return err
	}

	ctx.Export("sqsOperationQueuePolicy", sqsOperationQueuePolicy.ID())

	return nil
}
