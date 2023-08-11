package awsdeploy

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	SNSSignalsTopic      pulumi.StringOutput
	SNSPriceRequestTopic pulumi.StringOutput
	SNSOperationTopic    pulumi.StringOutput
)

func CreateTopics(ctx *pulumi.Context) error {
	snsSignalsTopic, err := sns.NewTopic(ctx, "SNSSignalsTopic", &sns.TopicArgs{
		Name:      pulumi.String("trading-bot-signals-topic.fifo"),
		FifoTopic: pulumi.Bool(true),
	})

	if err != nil {
		return err
	}

	SNSSignalsTopic = snsSignalsTopic.Arn
	ctx.Export("snsSignalsTopic", snsSignalsTopic.Arn)

	snsPriceRequestTopic, err := sns.NewTopic(ctx, "SNSPriceRequestTopic", &sns.TopicArgs{
		Name:      pulumi.String("trading-bot-price-request-topic.fifo"),
		FifoTopic: pulumi.Bool(true),
	})

	if err != nil {
		return err
	}

	SNSPriceRequestTopic = snsPriceRequestTopic.Arn
	ctx.Export("snsPriceRequestTopic", snsPriceRequestTopic.Arn)

	snsOperationTopic, err := sns.NewTopic(ctx, "SNSOperationTopic", &sns.TopicArgs{
		Name:      pulumi.String("trading-bot-operation-topic.fifo"),
		FifoTopic: pulumi.Bool(true),
	})

	if err != nil {
		return err
	}

	SNSOperationTopic = snsOperationTopic.Arn
	ctx.Export("snsOperationTopic", snsOperationTopic.Arn)

	return nil
}
