package pubsub

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sns"
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
