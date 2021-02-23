package health

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/julienbreux/rabdis/pkg/logger"
)

// Health represents the Health interface
type Health interface {
	Start()
	Stop() error
}

// health represents health
type health struct {
	o *Options

	httpServer    *http.Server
	ctxCancelFunc context.CancelFunc
}

// New creates a new Health instance
func New(opts ...Option) (Health, error) {
	o, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc(o.Route, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", o.Port),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	m := &health{
		o: o,

		httpServer: httpServer,

		ctxCancelFunc: cancel,
	}

	return m, nil
}

// Start starts health server
func (m *health) Start() {
	m.o.Logger.Info("health server: starting")

	go m.start()

	m.o.Logger.Info(
		"health server: started",
		logger.F("port", m.o.Port),
		logger.F("route", m.o.Route),
	) // Fake
}

func (m *health) start() {
	if err := m.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		m.o.Logger.Error(
			"health: unable to preserve health server started",
			logger.E(err),
		)
	}
}

// Stop stops health server
func (m *health) Stop() error {
	m.o.Logger.Info("health server: stopping")
	m.ctxCancelFunc()
	m.o.Logger.Info("health server: stopped")

	return nil
}
