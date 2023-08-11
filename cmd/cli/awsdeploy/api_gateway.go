package awsdeploy

import (
	"os"
	
	"jarvis-trading-bot/consts"

    "github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateApiGateway(ctx *pulumi.Context) error {
    
    account, err := aws.GetCallerIdentity(ctx)
	if err != nil {
		return err
	}

	region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
	if err != nil {
		return err
	}

    
	apiGatewayAccount, err := apigateway.NewAccount(ctx, "ApiGatewayAccountConfig", &apigateway.AccountArgs{
		CloudwatchRoleArn: ApiGatewayLoggingRole.Arn,
	})
	if err != nil {
		return err
	}

	apiGatewayRestApi, err := apigateway.NewRestApi(ctx, "ApiGatewayRestApi", &apigateway.RestApiArgs{
		Name: pulumi.String("trading-bot-service"),
		EndpointConfiguration: &apigateway.RestApiEndpointConfigurationArgs{
			Types: pulumi.String("REGIONAL"),
		},
		Policy: pulumi.String(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Action": "sts:AssumeRole",
                    "Principal": {
                        "Service": "lambda.amazonaws.com"
                    },
                    "Effect": "Allow",
                    "Sid": ""
                },
                {
                    "Action": "execute-api:Invoke",
                    "Resource": "*",
                    "Principal": "*",
                    "Effect": "Allow",
                    "Sid": ""
                }
            ]
        }`),
		MinimumCompressionSize: pulumi.Int(1024),
	},
		pulumi.DependsOn([]pulumi.Resource{apiGatewayAccount}),
	)
	if err != nil {
		return err
	}

	_, err = lambda.NewPermission(ctx, "AuthorizerHandlerLambdaApiGatewayPermission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  AuthorizerHandlerLambdaFunction.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, apiGatewayRestApi.ID()),
	})
	if err != nil {
		return err
	}

	_, err = lambda.NewPermission(ctx, "ExchangeConfigHandlerLambdaApiGatewayPermission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  ExchangeConfigHandlerLambdaFunction.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, apiGatewayRestApi.ID()),
	})
	if err != nil {
		return err
	}

	receiveSignalHandlerLambdaApiGatewayPermission, err := lambda.NewPermission(ctx, "ReceiveSignalHandlerLambdaApiGatewayPermission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  ReceiveSignalHandlerLambdaFunction.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, apiGatewayRestApi.ID()),
	})
	if err != nil {
		return err
	}

	receiveAlertSignalHandlerLambdaApiGatewayPermission, err := lambda.NewPermission(ctx, "ReceiveAlertSignalHandlerLambdaApiGatewayPermission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  ReceiveAlertSignalHandlerLambdaFunction.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, apiGatewayRestApi.ID()),
	})
	if err != nil {
		return err
	}

	apiGatewayAuthorizer, err := apigateway.NewAuthorizer(ctx, "ApiGatewayAuthorizer", &apigateway.AuthorizerArgs{
		AuthorizerResultTtlInSeconds: pulumi.Int(300),
		AuthorizerUri: pulumi.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations",
			pulumi.String(os.Getenv(consts.AwsRegion)),
			AuthorizerHandlerLambdaFunction.Arn,
		),
		IdentitySource: pulumi.String("method.request.header.Authorization"),
		Name:           pulumi.String("Authorizer"),
		RestApi:        apiGatewayRestApi.ID(),
		Type:           pulumi.String("REQUEST"),
	}, pulumi.DependsOn([]pulumi.Resource{AuthorizerHandlerLambdaFunction}))
	if err != nil {
		return err
	}

	apiGatewayResourceExchangeConfig, err := apigateway.NewResource(ctx, "ApiGatewayResourceExchangeConfig", &apigateway.ResourceArgs{
		ParentId: apiGatewayRestApi.RootResourceId,
		PathPart: pulumi.String("exchange-config"),
		RestApi:  apiGatewayRestApi.ID(),
	})
	if err != nil {
		return err
	}

	apiGatewayMethodExchangeConfigPost, err := apigateway.NewMethod(ctx, "ApiGatewayMethodExchangeConfigPost", &apigateway.MethodArgs{
		HttpMethod:     pulumi.String("POST"),
		ResourceId:     apiGatewayResourceExchangeConfig.ID(),
		RestApi:        apiGatewayRestApi.ID(),
		ApiKeyRequired: pulumi.Bool(false),
		Authorization:  pulumi.String("CUSTOM"),
		AuthorizerId:   apiGatewayAuthorizer.ID(),
	})
	if err != nil {
		return err
	}

	apiGatewayMethodExchangeConfigPostIntegration, err := apigateway.NewIntegration(ctx, "ApiGatewayMethodExchangeConfigPostIntegration", &apigateway.IntegrationArgs{
		HttpMethod:            apiGatewayMethodExchangeConfigPost.HttpMethod,
		IntegrationHttpMethod: pulumi.String("POST"),
		ResourceId:            apiGatewayResourceExchangeConfig.ID(),
		RestApi:               apiGatewayRestApi.ID(),
		Type:                  pulumi.String("AWS_PROXY"),
		Uri:                   ExchangeConfigHandlerLambdaFunction.InvokeArn,
	},
		pulumi.DependsOn([]pulumi.Resource{
			apiGatewayAuthorizer,
			apiGatewayMethodExchangeConfigPost,
			ExchangeConfigHandlerLambdaFunction,
		}),
	)
	if err != nil {
		return err
	}

	apiGatewayResourceSignal, err := apigateway.NewResource(ctx, "ApiGatewayResourceSignal", &apigateway.ResourceArgs{
		ParentId: apiGatewayRestApi.RootResourceId,
		PathPart: pulumi.String("signal"),
		RestApi:  apiGatewayRestApi.ID(),
	})
	if err != nil {
		return err
	}

	apiGatewayMethodSignalPost, err := apigateway.NewMethod(ctx, "ApiGatewayMethodSignalPost", &apigateway.MethodArgs{
		HttpMethod:     pulumi.String("POST"),
		ResourceId:     apiGatewayResourceSignal.ID(),
		RestApi:        apiGatewayRestApi.ID(),
		ApiKeyRequired: pulumi.Bool(false),
		Authorization:  pulumi.String("NONE"),
	}, pulumi.DependsOn([]pulumi.Resource{receiveSignalHandlerLambdaApiGatewayPermission}))
	if err != nil {
		return err
	}

	apiGatewayMethodSignalPostIntegration, err := apigateway.NewIntegration(ctx, "ApiGatewayMethodSignalPostIntegration", &apigateway.IntegrationArgs{
		HttpMethod:            apiGatewayMethodSignalPost.HttpMethod,
		IntegrationHttpMethod: pulumi.String("POST"),
		ResourceId:            apiGatewayResourceSignal.ID(),
		RestApi:               apiGatewayRestApi.ID(),
		Type:                  pulumi.String("AWS_PROXY"),
		Uri:                   ReceiveSignalHandlerLambdaFunction.InvokeArn,
	}, pulumi.DependsOn([]pulumi.Resource{
		apiGatewayMethodSignalPost,
		ReceiveSignalHandlerLambdaFunction,
	}))
	if err != nil {
		return err
	}

	apiGatewayResourcePostAlert, err := apigateway.NewResource(ctx, "ApiGatewayResourcePostAlert", &apigateway.ResourceArgs{
		ParentId: apiGatewayResourceSignal.ID(),
		PathPart: pulumi.String("{alert}"),
		RestApi:  apiGatewayRestApi.ID(),
	})
	if err != nil {
		return err
	}

	apiGatewayMethodAlertSignalPost, err := apigateway.NewMethod(ctx, "ApiGatewayMethodAlertSignalPost", &apigateway.MethodArgs{
		HttpMethod:     pulumi.String("POST"),
		ResourceId:     apiGatewayResourcePostAlert.ID(),
		RestApi:        apiGatewayRestApi.ID(),
		ApiKeyRequired: pulumi.Bool(false),
		Authorization:  pulumi.String("NONE"),
	}, pulumi.DependsOn([]pulumi.Resource{receiveAlertSignalHandlerLambdaApiGatewayPermission}))
	if err != nil {
		return err
	}

	apiGatewayMethodAlertSignalPostIntegration, err := apigateway.NewIntegration(ctx, "ApiGatewayMethodAlertSignalPostIntegration", &apigateway.IntegrationArgs{
		HttpMethod:            apiGatewayMethodAlertSignalPost.HttpMethod,
		IntegrationHttpMethod: pulumi.String("POST"),
		ResourceId:            apiGatewayResourcePostAlert.ID(),
		RestApi:               apiGatewayRestApi.ID(),
		Type:                  pulumi.String("AWS_PROXY"),
		Uri:                   ReceiveAlertSignalHandlerLambdaFunction.InvokeArn,
	}, pulumi.DependsOn([]pulumi.Resource{
		apiGatewayMethodAlertSignalPost,
		ReceiveAlertSignalHandlerLambdaFunction,
	}))
	if err != nil {
		return err
	}

	apiGatewayMethodDeployment, err := apigateway.NewDeployment(ctx, "ApiGatewayMethodDeployment", &apigateway.DeploymentArgs{
		Description: pulumi.String("Gateway API deployment"),
		RestApi:     apiGatewayRestApi.ID(),
	}, pulumi.DependsOn([]pulumi.Resource{
		apiGatewayMethodExchangeConfigPost,
		apiGatewayMethodExchangeConfigPostIntegration,
		apiGatewayMethodSignalPost,
		apiGatewayMethodSignalPostIntegration,
		apiGatewayMethodAlertSignalPost,
		apiGatewayMethodAlertSignalPostIntegration,
	}))
	if err != nil {
		return err
	}

	apiGatewayStage, err := apigateway.NewStage(ctx, "ApiGatewayStage", &apigateway.StageArgs{
		Deployment: apiGatewayMethodDeployment.ID(),
		RestApi:    apiGatewayRestApi.ID(),
		StageName:  pulumi.String("prod"),
		AccessLogSettings: &apigateway.StageAccessLogSettingsArgs{
			DestinationArn: ApiGatewayLogGroup.Arn,
			Format:         pulumi.String(`$context.extendedRequestId $context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] "$context.httpMethod $context.resourcePath $context.protocol" $context.status $context.responseLength $context.requestId`),
		},
	}, pulumi.DependsOn([]pulumi.Resource{apiGatewayMethodDeployment}))
	if err != nil {
		return err
	}

	apiGatewayMethodSettings, err := apigateway.NewMethodSettings(ctx, "ApiGatewayMethodSettings", &apigateway.MethodSettingsArgs{
		RestApi:    apiGatewayRestApi.ID(),
		StageName:  apiGatewayStage.StageName,
		MethodPath: pulumi.String("*/*"),
		Settings: &apigateway.MethodSettingsSettingsArgs{
			MetricsEnabled: pulumi.Bool(true),
			LoggingLevel:   pulumi.String("INFO"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{apiGatewayStage}))
	if err != nil {
		return err
	}

	apiDomainName, err := apigateway.NewDomainName(ctx, "ApiDomainName", &apigateway.DomainNameArgs{
		DomainName: pulumi.String("api.lumi.dev.br"),
		EndpointConfiguration: &apigateway.DomainNameEndpointConfigurationArgs{
			Types: pulumi.String("REGIONAL"),
		},
		RegionalCertificateArn: pulumi.String("arn:aws:acm:us-east-1:597887252921:certificate/a5b3521a-5354-4701-af28-63697f4c62e1"),
	}, pulumi.DependsOn([]pulumi.Resource{apiGatewayMethodSettings}))
	if err != nil {
		return err
	}
    ctx.Export("regionalDomainName", apiDomainName.RegionalDomainName)
    
    
	_, err = apigateway.NewBasePathMapping(ctx, "ApiBasePathMapping", &apigateway.BasePathMappingArgs{
		DomainName: pulumi.String("api.lumi.dev.br"),
		RestApi:    apiGatewayRestApi.ID(),
		BasePath:   pulumi.String("trading"),
		StageName:  pulumi.String("prod"),
	}, pulumi.DependsOn([]pulumi.Resource{apiDomainName}))
	if err != nil {
		return err
	}

	return nil
}
