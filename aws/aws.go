package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/hrsupersport/hrnogomet-backend-kit/constants"
)

// CustomEndpointResolver does custom endpoint resolution
// used to run aws calls against localstack
// see https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/endpoints
type CustomEndpointResolver struct {
	CustomEndpoint string
}

func (r *CustomEndpointResolver) ResolveEndpoint(serviceID, regionID string, options ...interface{}) (aws.Endpoint, error) {
	// You can implement custom logic here to return the endpoint based on the serviceID and regionID.
	return aws.Endpoint{
		URL:           r.CustomEndpoint,
		SigningRegion: regionID,
	}, nil
}

func newAwsConfig(ctx context.Context, awsRegion string, customEndpoint string) (*aws.Config, error) {
	var awsConfig aws.Config
	var err error

	if customEndpoint != "" {
		resolver := &CustomEndpointResolver{
			CustomEndpoint: customEndpoint,
		}
		awsConfig, err = config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion), config.WithEndpointResolverWithOptions(resolver))
	} else {
		awsConfig, err = config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	}

	if err != nil {
		return nil, err
	} else {
		return &awsConfig, nil
	}
}

// getCustomAwsEndpoint returns custom context key with custom aws endpoint if set
// use context.WithValue(ctx, constants.ContextKeyCustomAwsEndpoint, "<<endpoint>>") to set it
func getCustomAwsEndpoint(ctx context.Context) string {
	if val, ok := ctx.Value(constants.ContextKeyCustomAwsEndpoint{}).(string); ok {
		return val
	}
	return ""
}

// SetCustomAwsEndpoint sets custom context attribute for aws custom endpoint
func SetCustomAwsEndpoint(ctx context.Context, customEndpoint string) context.Context {
	return context.WithValue(ctx, constants.ContextKeyCustomAwsEndpoint{}, customEndpoint)
}

// CreateDynamodbClient creates new dynamodb client
func CreateDynamodbClient(ctx context.Context, awsRegion string) (*dynamodb.Client, error) {
	if awsConfig, err := newAwsConfig(ctx, awsRegion, getCustomAwsEndpoint(ctx)); err != nil {
		return nil, err
	} else {
		return dynamodb.NewFromConfig(*awsConfig), nil
	}
}

// CreateSqsClient creates new AWS SQS client
func CreateSqsClient(ctx context.Context, awsRegion string) (*sqs.Client, error) {
	if awsConfig, err := newAwsConfig(ctx, awsRegion, getCustomAwsEndpoint(ctx)); err != nil {
		return nil, err
	} else {
		return sqs.NewFromConfig(*awsConfig), nil
	}
}
