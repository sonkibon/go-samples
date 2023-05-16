package pubsub

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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
