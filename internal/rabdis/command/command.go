package command

import (
	"io"
	"log"

	"github.com/julienbreux/rabdis/internal/rabdis"
	"github.com/julienbreux/rabdis/internal/rabdis/command/version"
	"github.com/julienbreux/rabdis/pkg/logger"
	"github.com/julienbreux/rabdis/pkg/metrics"
	"github.com/julienbreux/rabdis/pkg/rabbitmq"
	"github.com/julienbreux/rabdis/pkg/redis"
	ver "github.com/julienbreux/rabdis/pkg/version"
	"github.com/spf13/cobra"
)

// NewCmdRoot returns the Hut CLI
func NewCmdRoot(in io.Reader, out, err io.Writer) (cmd *cobra.Command) {
	// Create CLI
	cmd = &cobra.Command{
		Use:     "rabdis",
		Short:   "rabdis is a program that help to delete Redis keys from RabbitMQ messages",
		Version: ver.Version,
		Run:     RunCmdRoot,
	}

	cmd.AddCommand(version.NewCmdVersion(in, out, err))

	return
}

// RunCmdRoot runs the root command
func RunCmdRoot(cmd *cobra.Command, args []string) {
	// Logger
	lgr, err := logger.New(
		logger.InstName(cmd.Use),
		logger.InstVersion(ver.Version),
	)
	if err != nil {
		log.Fatal("rabdis: logger configuration failed", logger.E(err))
	}

	// Manager
	rbds, err := rabdis.New(
		rabdis.Logger(lgr),
	)
	if err != nil {
		log.Fatal("rabdis: configuration failed", logger.E(err))
	}

	// RabbitMQ
	rmq, err := rabbitmq.New(
		rabbitmq.Logger(lgr),
		rabbitmq.InstName(cmd.Use),
		rabbitmq.InstVersion(cmd.Version),
	)
	if err != nil {
		log.Fatal("rabbitmq: configuration failed", logger.E(err))
	}
	rbds.SetRabbitMQ(rmq)

	// Redis
	red, err := redis.New(
		redis.Logger(lgr),
	)
	if err != nil {
		log.Fatal("redis: configuration failed", logger.E(err))
	}
	rbds.SetRedis(red)

	// Metrics
	met, err := metrics.New(
		metrics.Logger(lgr),
	)
	if err != nil {
		log.Fatal("metrics: configuration failed", logger.E(err))
	}
	rbds.SetMetrics(met)

	// Let's go baby!
	rbds.Start()
}
