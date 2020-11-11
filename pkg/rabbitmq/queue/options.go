package queue

import "context"

// DurableQueueKey represents a durable queue key used by context
type DurableQueueKey struct{}

// DeadletterQueueKey represents a deadletter queue key used by context
type DeadletterQueueKey struct{}

// MessagesTTLQueueKey represents a message TTL queue key used by context
type MessagesTTLQueueKey struct{}

// QueueOptions represents the queue options
type QueueOptions struct {
	Context context.Context
}

// QueueOption represents an option
type QueueOption func(*QueueOptions)

// QueueDurable creates a durable queue
func QueueDurable() QueueOption {
	return func(o *QueueOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, DurableQueueKey{}, true)
	}
}

// QueueDeadletter creates a deadletter queue
func QueueDeadletter(exchange string) QueueOption {
	return func(o *QueueOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, DeadletterQueueKey{}, exchange)
	}
}

// QueueMessagesTTL creates a TTL on messages
func QueueMessagesTTL(ttl int32) QueueOption {
	return func(o *QueueOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, MessagesTTLQueueKey{}, ttl)
	}
}
