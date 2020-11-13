package signal

import (
	"os"
	"os/signal"
)

const sigChanSize = 2

var singleton = make(chan struct{})

// Handler returns channel SIGTERM and SIGINT.
// If a second signal is caught, terminate program with exit code 1.
func Handler() <-chan struct{} {
	close(singleton) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, sigChanSize)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}
