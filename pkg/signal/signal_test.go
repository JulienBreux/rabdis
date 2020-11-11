package signal_test

import (
	"testing"

	"github.com/julienbreux/rabdis/pkg/signal"
	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	signal.Handler()
	assert.Panics(t, func() { signal.Handler() })
}
