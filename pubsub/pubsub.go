package pubsub

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	NameQueueArn                = "QueueArn"
	NameTopicArn                = "TopicArn"
	QueueAttributeRedrivePolicy = "RedrivePolicy"
)

var (
	SubscriptionProtocolSQS   = aws.String("sqs")
	SubscriptionProtocolHTTP  = aws.String("http")
	SubscriptionProtocolHTTPS = aws.String("https")
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

// NewQueue calls the NewTopicContext method.
func (c *PubsubClient) NewTopic(topicArn string) (*Topic, error) {
	return c.NewTopicContext(context.Background(), topicArn)
}

// NewTopicContext returns an initialized topic client based on the topic arn.
func (c *PubsubClient) NewTopicContext(ctx context.Context, topicArn string) (*Topic, error) {
	parse, err := arn.Parse(topicArn)
	if err != nil {
		return nil, fmt.Errorf("arn.Parse(%s) : %w", topicArn, err)
	}

	if _, err = c.SNS.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{TopicArn: &topicArn}); err != nil {
		return nil, fmt.Errorf("c.SNS.GetTopicAttributes: %w", err)
	}

	return &Topic{
		client:    c,
		topicName: parse.Resource,
		topicArn:  topicArn,
	}, nil
}

// NewSubscription calls the NewSubscriptionContext method.
func (c *PubsubClient) NewSubscription(subscriptionArn string) (*Subscription, error) {
	return c.NewSubscriptionContext(context.Background(), subscriptionArn)
}

// NewSubscriptionContext returns an initialized subscription client based on the subscription arn.
func (c *PubsubClient) NewSubscriptionContext(ctx context.Context, subscriptionArn string) (*Subscription, error) {
	if _, err := arn.Parse(subscriptionArn); err != nil {
		return nil, fmt.Errorf("arn.Parse: %w", err)
	}

	atr, err := c.SNS.GetSubscriptionAttributes(
		ctx,
		&sns.GetSubscriptionAttributesInput{
			SubscriptionArn: &subscriptionArn,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("c.SNS.GetSubscriptionAttributes(%s) : %w", subscriptionArn, err)
	}
	topic, err := c.NewTopic(atr.Attributes[NameTopicArn])
	if err != nil {
		return nil, fmt.Errorf("c.NewTopic: %w", err)
	}

	subscriptions, err := c.SNS.ListSubscriptionsByTopic(
		ctx,
		&sns.ListSubscriptionsByTopicInput{
			TopicArn: &topic.topicArn,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("c.SNS.ListSubscriptionsByTopic(%s) : %w", topic.topicArn, err)
	}

	for _, subscription := range subscriptions.Subscriptions {
		if *subscription.Protocol == *SubscriptionProtocolSQS && *subscription.SubscriptionArn == subscriptionArn {
			queue, err := c.NewQueue(*subscription.Endpoint)
			if err != nil {
				return nil, fmt.Errorf("c.NewQueue(%v) : %w", subscription.Endpoint, err)
			}

			return &Subscription{
				client:          c,
				subscriptionArn: subscriptionArn,
				topic:           *topic,
				queue:           *queue,
			}, nil
		}
	}

	return nil, errors.New("subscription not found")
}

// Change old opts format to new format
func (c *PubsubClient) convertOldOpts(in map[string]*string) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		if v == nil {
			out[k] = ""
		} else {
			out[k] = *v
		}
	}

	return out
}

// CreateQueue calls the CreateQueueContext method.
func (c *PubsubClient) CreateQueue(queueName string, opts map[string]*string) (*Queue, error) {
	return c.CreateQueueContext(context.Background(), queueName, opts)
}

// CreateQueueContext returns an initialized queue client based on the queue name and options.
func (c *PubsubClient) CreateQueueContext(ctx context.Context, queueName string, opts map[string]*string) (*Queue, error) {
	queue, err := c.SQS.CreateQueue(
		ctx,
		&sqs.CreateQueueInput{
			QueueName:  &queueName,
			Attributes: c.convertOldOpts(opts),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("c.SQS.CreateQueue(queueName=%s, attributes=%+v) : %w", queueName, opts, err)
	}

	atr, err := c.SQS.GetQueueAttributes(
		ctx,
		&sqs.GetQueueAttributesInput{
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
			QueueUrl:       queue.QueueUrl,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("c.SQS.GetQueueAttributes(queueUrl=%v) : %w", queue.QueueUrl, err)
	}

	return &Queue{
		client:    c,
		queueArn:  atr.Attributes[NameQueueArn],
		queueName: queueName,
		queueUrl:  *queue.QueueUrl,
	}, nil
}
