package metrics

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics represents the Metrics interface
type Metrics interface {
	Start()
	Stop() error
}

// metrics represents metrics
type metrics struct {
	o *Options

	httpServer    *http.Server
	ctxCancelFunc context.CancelFunc
}

// New creates a new Metrics instance
func New(opts ...Option) (Metrics, error) {
	o, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	mux.Handle(o.Route, promhttp.Handler())
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", o.Port),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	m := &metrics{
		o: o,

		httpServer: httpServer,

		ctxCancelFunc: cancel,
	}

	return m, nil
}

// Start starts metrics server
func (m *metrics) Start() {
	m.o.Logger.Info("metrics server: starting")

	go m.start()

	m.o.Logger.Info(
		"metrics server: started",
		logger.F("port", m.o.Port),
		logger.F("route", m.o.Route),
	) // Fake
}

func (m *metrics) start() {
	if err := m.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		m.o.Logger.Error(
			"metrics: unable to preserve metrics server started",
			logger.E(err),
		)
	}
}

// Stop stops metrics server
func (m *metrics) Stop() error {
	m.o.Logger.Info("metrics server: stopping")
	m.ctxCancelFunc()
	m.o.Logger.Info("metrics server: stopped")

	return nil
}
