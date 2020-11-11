package url_test

import (
	"testing"

	"github.com/julienbreux/rabdis/pkg/url"
	"github.com/stretchr/testify/assert"
)

func TestValuesToURLAll(t *testing.T) {
	actual := url.Build("mongodb", "julien", "breux", "0.0.0.0", 27017, nil)
	expected := "mongodb://julien:breux@0.0.0.0:27017"

	assert.Equal(t, expected, actual)
}

func TestValuesToURLNoPort(t *testing.T) {
	actual := url.Build("mongodb", "julien", "breux", "0.0.0.0", 0, nil)
	expected := "mongodb://julien:breux@0.0.0.0"

	assert.Equal(t, expected, actual)
}

func TestValuesToURLNoPassword(t *testing.T) {
	actual := url.Build("mongodb", "julien", "", "0.0.0.0", 27017, nil)
	expected := "mongodb://julien@0.0.0.0:27017"

	assert.Equal(t, expected, actual)
}

func TestValuesToURLNoUsername(t *testing.T) {
	actual := url.Build("mongodb", "", "breux", "0.0.0.0", 27017, nil)
	expected := "mongodb://0.0.0.0:27017"

	assert.Equal(t, expected, actual)
}

func TestValuesToUrlPath(t *testing.T) {
	path := "/"
	actual := url.Build("mongodb", "", "breux", "0.0.0.0", 27017, &path)
	expected := "mongodb://0.0.0.0:27017" + path

	assert.Equal(t, expected, actual)
}
