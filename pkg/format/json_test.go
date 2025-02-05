package format_test

import (
	"testing"

	"github.com/julienbreux/rabdis/pkg/format"
	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	data := struct {
		Version string `json:"version"`
	}{
		Version: "1.0.0",
	}
	actual, err := format.ToJSON(data)
	expected := "{\"version\":\"1.0.0\"}"

	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}
