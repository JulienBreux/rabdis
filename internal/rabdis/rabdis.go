package rabdis

import (
	"sync"

	"github.com/julienbreux/rabdis/internal/rabdis/config"
	"github.com/julienbreux/rabdis/pkg/health"
	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/julienbreux/rabdis/pkg/metrics"
	"github.com/julienbreux/rabdis/pkg/rabbitmq"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/consumer"
	"github.com/julienbreux/rabdis/pkg/redis"
	"github.com/julienbreux/rabdis/pkg/signal"
	"github.com/julienbreux/rabdis/pkg/version"
)

const numServices = 3

// Rabdis represents the Rabdis interface
type Rabdis interface {
	SetRabbitMQ(rabbitmq.RabbitMQ)
	SetRedis(redis.Redis)
	SetMetrics(metrics.Metrics)
	SetHealth(health.Health)

	Start()
}

type rabdis struct {
	o *Options
	c *config.Config

	doneCh chan struct{}
	stopCh <-chan struct{}

	wg sync.WaitGroup

	rabbitMQ rabbitmq.RabbitMQ
	redis    redis.Redis
	metrics  metrics.Metrics
	health   health.Health
}

// New creates a new Rabdis instance
func New(opts ...Option) (Rabdis, error) {
	o, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	c, err := config.FromFile(o.ConfigFile, o.Logger)
	if err != nil {
		return nil, err
	}

	return &rabdis{
		o: o,
		c: c,

		doneCh: make(chan struct{}),
		stopCh: signal.Handler(),
	}, nil
}

// SetRabbitMQ sets RabbitMQ
func (r *rabdis) SetRabbitMQ(rmq rabbitmq.RabbitMQ) {
	r.rabbitMQ = rmq
}

// SetRedis sets Redis
func (r *rabdis) SetRedis(rds redis.Redis) {
	r.redis = rds
}

// SetMetrics sets Metrics
func (r *rabdis) SetMetrics(mts metrics.Metrics) {
	r.metrics = mts
}

// SetHealth sets Health
func (r *rabdis) SetHealth(hlth health.Health) {
	r.health = hlth
}

// Start starts Rabdis
func (r *rabdis) Start() {
	r.o.Logger.Info(
		"rabdis: starting",
		logger.F("version", version.Version),
		logger.F("build", version.Commit),
	)

	// TODO: manage RabbitMQ hook error
	r.init()
	r.start()

	r.o.Logger.Info("rabdis: started")

	select {
	// Signal stop
	case <-r.stopCh:
		r.o.Logger.Debug("rabdis: stop by signal")
		break
	// Rabdis properly done
	case <-r.doneCh:
		r.o.Logger.Debug("rabdis: stop by done")
		break
	}

	r.stop()
}

func (r *rabdis) init() {
	for _, rule := range r.c.Rules {
		r.rabbitMQ.OnMessage(
			r.messageHandler(rule),
			consumer.Exchange(rule.RabbitMQ.Exchange),
			consumer.Queue(rule.RabbitMQ.Queue),
			consumer.Bind(rule.RabbitMQ.Bind),
		)
	}
}

func (r *rabdis) start() {
	r.wg.Add(numServices)

	go r.rabbitMQ.Connect()
	go r.redis.Connect()
	go r.metrics.Start()
	go r.health.Start()

	go func() {
		r.wg.Wait()
		close(r.doneCh)
	}()
}

// stop stops Rabdis
func (r *rabdis) stop() {
	r.o.Logger.Info("rabdis: stopping")

	// TODO: manager disconnection errors
	_ = r.rabbitMQ.Disconnect()
	_ = r.redis.Disconnect()
	_ = r.metrics.Stop()
	_ = r.health.Stop()

	r.o.Logger.Info("rabdis: stopped")
}
