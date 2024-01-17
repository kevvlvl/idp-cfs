package global

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogError(t *testing.T) {

	e := LogError("this is a test error")
	assert.NotNil(t, e)
	assert.Contains(t, e.Error(), "this is a test error")
}
