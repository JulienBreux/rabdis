package rabbitmq

import (
	"net"
	"time"

	"github.com/jpillora/backoff"
	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/consumer"
	"github.com/julienbreux/rabdis/pkg/rabbitmq/message"
	"github.com/julienbreux/rabdis/pkg/url"
	"github.com/streadway/amqp"
)

// RabbitMQ represents the RabbitMQ interface
type RabbitMQ interface {
	Connect()
	Disconnect() error

	OnMessage(message.OnMessageHandler, ...consumer.Option) error
}

// OnConnectHandler represents a handler when RabbitMQ is connected
type OnConnectHandler func(RabbitMQ) error

// rabbitmq represents the internal rabbitmq structure
type rabbitmq struct {
	o *Options

	conn *amqp.Connection

	consumers []consumer.Consumer
}

// New creates a new RabbitMQ instance
func New(opts ...Option) (RabbitMQ, error) {
	o, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	r := &rabbitmq{
		o: o,
	}

	return r, nil
}

// Connect connects to the RabbitMQ server
func (r *rabbitmq) Connect() {
	r.connect()
}

func (r *rabbitmq) connect() {
	r.o.Logger.Info("rabbitmq: connecting")

	b := &backoff.Backoff{}
	attempts := 1
	for {
		var err error
		url := url.Build("amqp", r.o.Username, r.o.Password, r.o.Host, r.o.Port, &r.o.VirtualHost)
		conn, err := amqp.DialConfig(url, r.config())
		if err == nil {
			r.conn = conn
			r.o.Logger.Info("rabbitmq: connected")
			for _, c := range r.consumers {
				_ = c.Start(r.conn)
			}
			break
		}

		r.o.Logger.Error(
			"rabbitmq: failed to connect",
			logger.F("attempts", attempts),
			logger.E(err),
		)

		attempts++
		time.Sleep(b.Duration())
	}
	r.watchConnectCloseAndReconnect()
}

func (r *rabbitmq) watchConnectCloseAndReconnect() {
	go func() {
		reason, ok := <-r.conn.NotifyClose(make(chan *amqp.Error))
		r.conn = nil
		if !ok {
			r.o.Logger.Debug("rabbitmq: connection close properly")
			return
		}
		r.o.Logger.Warn(
			"rabbitmq: disconnected",
			logger.F("reason", reason),
		)

		r.connect()
	}()
}

// config returns RabbitMQ configuration
func (r *rabbitmq) config() amqp.Config {
	connTimeout := time.Duration(r.o.ConnTimeout) * time.Second
	return amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, connTimeout)
		},
		Properties: amqp.Table{
			"product": r.o.InstName,
			"version": r.o.InstVersion,
		},
	}
}

// Disconnect disconnects from RabbitMQ
func (r *rabbitmq) Disconnect() error {
	if r.conn == nil {
		return nil
	}

	r.o.Logger.Info("rabbitmq: disconnecting")

	err := r.conn.Close()
	if err == nil {
		r.o.Logger.Info("rabbitmq: disconnected")
		return nil
	}

	r.o.Logger.Warn(
		"rabbitmq: unable to stop",
		logger.E(err),
	)

	return err
}

// OnMessage stores the on message receive handlers
// TODO: Add error management
func (r *rabbitmq) OnMessage(h message.OnMessageHandler, opts ...consumer.Option) error {
	c, _ := consumer.New(h, opts...)
	r.consumers = append(r.consumers, c)
	return nil
}
