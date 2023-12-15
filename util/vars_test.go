package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoolPtr(t *testing.T) {

	b := BoolPtr(true)
	assert.Equal(t, true, *b)
}

func TestStringPtr(t *testing.T) {

	s := StringPtr("Test String")
	assert.Equal(t, "Test String", *s)
}
