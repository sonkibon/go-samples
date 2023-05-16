package pubsub

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
