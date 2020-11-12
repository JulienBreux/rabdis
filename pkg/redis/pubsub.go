package redis

import "github.com/go-redis/redis"

// PubSubPublish helps to publish
func (r *red) PubSubPublish(channel string, message string) error {
	if r.r == nil {
		return errorConnection
	}

	pubsub := r.r.Subscribe(channel)

	// Check subscription
	if _, err := pubsub.Receive(); err != nil {
		return err
	}

	if err := r.r.Publish(channel, message).Err(); err != nil {
		return err
	}

	return nil
}

// PubSubSubscribe helps to subscribe
func (r *red) PubSubSubscribe(channel string, message chan<- string) error {
	if r.r == nil {
		return errorConnection
	}

	pubsub := r.r.Subscribe(channel)

	// Check subscription
	if _, err := pubsub.Receive(); err != nil {
		return err
	}

	go r.pubSub(pubsub.Channel(), message)

	return nil
}

func (r *red) pubSub(msgChan <-chan *redis.Message, message chan<- string) {
	for {
		if msg, ok := <-msgChan; ok {
			message <- msg.Payload
			break
		}
		r.o.Logger.Debug("PubSubSubscribe reading failed")
	}
}
