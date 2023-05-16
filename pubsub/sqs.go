package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"golang.org/x/sync/errgroup"
)

// Queue provides a PubsubClient for a specific queue.
type Queue struct {
	client    *PubsubClient
	queueArn  string
	queueName string
	queueUrl  string
}

// SNSEvent is the struct to map when sending messages to the queue via topic.
type SNSEvent struct {
	Type              string
	MessageId         string
	Message           string
	Token             string
	TopicArn          string
	SubscribeURL      *string
	Timestamp         string
	Signature         string
	SignatureVersion  string
	SigningCertURL    string
	MessageAttributes map[string]map[string]string
}

// S3Event is the struct to map when sending messages to the queue via s3.
type S3Event struct {
	Records []struct {
		AwsRegion         string                 `json:"awsRegion"`
		EventName         string                 `json:"eventName"`
		EventTime         string                 `json:"eventTime"`
		EventSource       string                 `json:"eventSource"`
		EventVersion      string                 `json:"eventVersion"`
		UserIdentity      map[string]string      `json:"userIdentity"`
		RequestParameters map[string]interface{} `json:"requestParameters"`
		ResponseElements  map[string]interface{} `json:"responseElements"`
		S3                struct {
			S3SchemaVersion string `json:"s3SchemaVersion"`
			ConfigurationId string `json:"configurationId"`
			Bucket          struct {
				Name          string                 `json:"name"`
				OwnerIdentity map[string]interface{} `json:"ownerIdentity"`
				ARN           string                 `json:"arn"`
			} `json:"bucket"`
			Object struct {
				Key       string `json:"key"`
				Size      int64  `json:"size"`
				ETag      string `json:"etag"`
				Sequencer string `json:"sequencer"`
			} `json:"object"`
		} `json:"s3"`
	}
}

// Exist returns whether the topic exists or not.
func (q *Queue) Exist(ctx context.Context) (bool, error) {
	if _, err := q.client.SQS.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(q.queueUrl),
	}); err != nil {
		return false, fmt.Errorf("q.client.SQS.GetQueueAttributes: %w", err)
	}

	return true, nil
}

// Send delivers a message to the specified queue.
func (q *Queue) Send(ctx context.Context, message string, attributes map[string]types.MessageAttributeValue) error {
	m, err := q.client.SQS.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:       aws.String(message),
		MessageAttributes: attributes,
		QueueUrl:          &q.queueUrl,
	})
	if err != nil {
		return fmt.Errorf("q.client.SQS.SendMessage: %w", err)
	}

	log.Default().Printf("message id: %s", *m.MessageId)
	return nil
}

// Consume calls the consume method.
func (q *Queue) Consume(ctx context.Context, handler func(c context.Context, message string) (bool, error)) error {
	return q.consume(ctx, func(ctx context.Context, m types.Message) (bool, error) {
		return handler(ctx, *m.Body)
	})
}

// ConsumeViaSNS maps the message to an SNSEvent struct and calls the consume method.
func (q *Queue) ConsumeViaSNS(ctx context.Context, handler func(c context.Context, event SNSEvent) (retryable bool, err error)) error {
	return q.consume(ctx, func(ctx context.Context, m types.Message) (bool, error) {
		var event SNSEvent
		err := json.Unmarshal([]byte(*m.Body), &event)
		if err != nil {
			log.Default().Printf("failed to unmarshal json, body: %s", *m.Body)
			return true, fmt.Errorf("json.Unmarshal: %w", err)
		}
		return handler(ctx, event)
	})
}

// ConsumeViaS3 maps the message to an S3Event struct and calls the consume method.
func (q *Queue) ConsumeViaS3(ctx context.Context, handler func(c context.Context, event S3Event) (retryable bool, err error)) error {
	return q.consume(ctx, func(ctx context.Context, m types.Message) (bool, error) {
		var event S3Event
		err := json.Unmarshal([]byte(*m.Body), &event)
		if err != nil {
			log.Default().Printf("failed to unmarshal json, body: %s", *m.Body)
			return true, fmt.Errorf("json.Unmarshal: %w", err)
		}
		return handler(ctx, event)
	})
}

// consume receives a message from a specific queue and executes the argument f function to delete the message.
// It can also retry by changing the visibility timeout of the specified message in the queue to a new value.
func (q *Queue) consume(ctx context.Context, f func(context.Context, types.Message) (bool, error)) error {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.queueUrl),
		MaxNumberOfMessages: q.client.Config.MaxNumberOfMessages,
		WaitTimeSeconds:     q.client.Config.WaitTimeSeconds,
	}

	output, err := q.client.SQS.ReceiveMessage(ctx, params)
	if err != nil {
		return fmt.Errorf("q.SQS.ReceiveMessage: %w", err)
	}

	group, _ := errgroup.WithContext(ctx)
	for _, message := range output.Messages {
		group.Go(func(m types.Message) func() error {
			return func() error {
				retryable, err := f(ctx, m)

				if err == nil {
					if _, err := q.client.SQS.DeleteMessage(ctx, &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(q.queueUrl),
						ReceiptHandle: m.ReceiptHandle,
					}); err != nil {
						return fmt.Errorf("q.SQS.DeleteMessage: %w", err)
					}

					return nil
				}

				if retryable {
					if _, err := q.client.SQS.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
						QueueUrl:          aws.String(q.queueUrl),
						ReceiptHandle:     m.ReceiptHandle,
						VisibilityTimeout: q.client.Config.RequeueVisibilityTimeout,
					}); err != nil {
						return fmt.Errorf("q.SQS.ChangeMessageVisibility: %w", err)
					}
				} else {
					if _, err := q.client.SQS.DeleteMessage(ctx, &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(q.queueUrl),
						ReceiptHandle: m.ReceiptHandle,
					}); err != nil {
						return fmt.Errorf("q.SQS.DeleteMessage: %w", err)
					}
				}
				return nil
			}
		}(message))
	}
	if err := group.Wait(); err != nil {
		return fmt.Errorf("group.Wait: %w", err)
	}

	return nil
}
