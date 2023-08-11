package awsdeploy

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	IamRoleLambdaExecution   pulumi.StringOutput
	IamPolicyLambdaExecution *iam.RolePolicy

	ApiGatewayLoggingRole *iam.Role
)

func CreateRolesPolicies(ctx *pulumi.Context) error {
	err := CreateLambdaRolePolicy(ctx)
	if err != nil {
		return err
	}

	err = CreateApiGatewayRolePolicy(ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateLambdaRolePolicy(ctx *pulumi.Context) error {
	role, err := iam.NewRole(ctx, "IamRoleLambdaExecution", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Sid": "",
				"Effect": "Allow",
				"Principal": {
					"Service": "lambda.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),
	})
	if err != nil {
		return err
	}
	IamRoleLambdaExecution = role.Arn

	lambdaPolicy, err := iam.NewRolePolicy(ctx, "IamPolicyLambdaExecution", &iam.RolePolicyArgs{
		Role: role.Name,
		Policy: pulumi.String(`{
							"Version": "2012-10-17",
							"Statement": [{
									"Effect": "Allow",
									"Action": [
											"logs:CreateLogGroup",
											"logs:CreateLogStream",
											"logs:PutLogEvents"
									],
									"Resource": "arn:aws:logs:*:*:*"
							},
							{
								"Effect": "Allow",
								"Action": [
										"s3:DeleteObject",
										"s3:GetObject",
										"s3:PutObject"
								],
								"Resource": "arn:aws:s3:::*"
							},
							{
								"Effect": "Allow",
								"Action": [
                  	"dynamodb:PutItem",
                  	"dynamodb:UpdateItem",
                  	"dynamodb:DeleteItem",
                  	"dynamodb:BatchWriteItem",
                  	"dynamodb:GetItem",
                  	"dynamodb:BatchGetItem",
                  	"dynamodb:Scan",
                  	"dynamodb:Query"								
								],
								"Resource": "arn:aws:dynamodb:*:*:*"
							},
							{
								"Effect": "Allow",
								"Action": [
										"sns:Publish"
								],
								"Resource": "arn:aws:sns:*:*:*"
							},
							{
								"Effect": "Allow",
								"Action": [
										"sqs:*"
								],
								"Resource": "arn:aws:sqs:*:*:*"
							}]
					}`),
	})

	if err != nil {
		return err
	}
	IamPolicyLambdaExecution = lambdaPolicy

	return nil
}

func CreateApiGatewayRolePolicy(ctx *pulumi.Context) error {
	role, err := iam.NewRole(ctx, "ApiGatewayLoggingRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Sid": "",
				"Effect": "Allow",
				"Principal": {
					"Service": "apigateway.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"),
		},
	})
	if err != nil {
		return err
	}
	ApiGatewayLoggingRole = role

	return nil
}
