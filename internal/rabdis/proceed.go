package rabdis

import (
	"github.com/julienbreux/rabdis/internal/rabdis/config"
	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/message"
)

// TODO: replace ruleIndex by using context
func (r *rabdis) messageHandler(rule config.Rule) message.OnMessageHandler {
	return func(msg message.Message) error {
		r.o.Logger.Debug(
			"rabdis: OnMessageHandler triggered",
			logger.F("message", msg),
		)

		for _, a := range rule.Redis.Actions {
			// Check condition
			a.SetContent(msg.Body().Raw())
			if !a.ConditionsCheck() {
				r.o.Logger.Warn(
					"rabdis: action skip because not match conditions",
					logger.F("key", a.Key),
					logger.F("action", a.Action),
				)
				continue
			}

			key, err := a.FinalKey()
			if err != nil {
				r.o.Logger.Error(
					"rabdis: action skip because variable used in key not found in message",
					logger.F("key", a.Key),
					logger.F("message", msg.Body().String()),
					logger.E(err),
				)
				continue
			}
			_, _ = r.action(a.Action, key)
		}
		_ = msg.Ack()

		return nil
	}
}

// action proceeds to Redis action
// TODO: add error context
func (r *rabdis) action(a config.ActionRedis, key string) (res int64, err error) {
	switch a {
	case config.ActionDelete:
		res, err = r.redis.Del(key)
	case config.ActionIncrement:
		res, err = r.redis.Increment(key)
	case config.ActionDecrement:
		res, err = r.redis.Decrement(key)
	}

	return
}
