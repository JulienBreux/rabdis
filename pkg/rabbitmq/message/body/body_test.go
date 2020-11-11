package body_test

import (
	"testing"

	"github.com/julienbreux/rabdis/pkg/rabbitmq/message/body"
	"github.com/stretchr/testify/assert"
)

func TestBody(t *testing.T) {
	r := []byte("message")
	b := body.New(r)

	assert.Equal(t, r, b.Raw())
	assert.Equal(t, "message", b.String())
}
