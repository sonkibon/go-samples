package pubsub

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

// Topic provides a PubsubClient for a specific topic.
type Topic struct {
	client    *PubsubClient
	topicName string
	topicArn  string
}

// Subscription provides a PubsubClient for a specific subscription.
type Subscription struct {
	client          *PubsubClient
	subscriptionArn string
	topic           Topic
	queue           Queue
}

// Exist returns whether the topic exists or not.
func (t *Topic) Exist(ctx context.Context) (bool, error) {
	_, err := t.client.SNS.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: &t.topicArn,
	})
	if err != nil {
		return false, fmt.Errorf("t.client.SNS.GetTopicAttributes: %w", err)
	}

	return true, nil
}

// Publish sends a message to an Amazon SNS topic, a text message.
func (t *Topic) Publish(ctx context.Context, message string, attributes map[string]types.MessageAttributeValue) error {
	m, err := t.client.SNS.Publish(ctx, &sns.PublishInput{
		Message:           aws.String(message),
		MessageAttributes: attributes,
		TopicArn:          &t.topicArn,
	})
	if err != nil {
		return fmt.Errorf("t.client.SNS.Publish: %w", err)
	}

	log.Default().Printf("message id: %s", *m.MessageId)
	return nil
}
