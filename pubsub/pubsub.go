package pubsub

import (
	"time"

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
