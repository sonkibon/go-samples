package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// sqsConfig represents sqs config.
type sqsConfig struct {
	MaxNumberOfMessages      int32
	WaitTimeSeconds          int32
	RequeueVisibilityTimeout int32
	MaxReceiveCount          int32
}

// PubsubClient provides the API clients to make operations call for Amazon Simple
// Queue Service and Amazon Simple Notification Service
type PubsubClient struct {
	SQS        *sqs.Client
	SNS        *sns.Client
	Config     sqsConfig
	OpsTimeout time.Duration
}

// NewPubsubClient returns a new client from the provided clients and config.
func NewPubsubClient(sqs *sqs.Client, sns *sns.Client, cfg sqsConfig) (*PubsubClient, error) {
	return &PubsubClient{
		SQS:    sqs,
		SNS:    sns,
		Config: cfg,
	}, nil
}

// NewQueue calls the NewQueueContext method.
func (c *PubsubClient) NewQueue(queueArn string) (*Queue, error) {
	return c.NewQueueContext(context.Background(), queueArn)
}

// NewQueueContext returns an initialized queue client based on the queue arn.
func (c *PubsubClient) NewQueueContext(ctx context.Context, queueArn string) (*Queue, error) {
	parse, err := arn.Parse(queueArn)
	if err != nil {
		return nil, fmt.Errorf("arn.Parse: %w", err)
	}

	queueUrl, err := c.SQS.GetQueueUrl(
		ctx,
		&sqs.GetQueueUrlInput{
			QueueName:              &parse.Resource,
			QueueOwnerAWSAccountId: &parse.AccountID,
		},
	)
	if err != nil {
		return nil,
			fmt.Errorf("c.SQS.GetQueueUrl: %w", err)
	}

	return &Queue{
		client:    c,
		queueName: parse.Resource,
		queueUrl:  *queueUrl.QueueUrl,
		queueArn:  queueArn,
	}, nil
}
