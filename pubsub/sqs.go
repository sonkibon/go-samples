package pubsub

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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
