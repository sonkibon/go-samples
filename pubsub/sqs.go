package pubsub

import (
	"context"
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
