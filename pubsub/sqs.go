package pubsub

// Queue provides a PubsubClient for a specific queue.
type Queue struct {
	client    *PubsubClient
	queueArn  string
	queueName string
	queueUrl  string
}
