// File: condition_test.go
// Tests for Condition
package gorule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextSetGetValue(t *testing.T) {
	ctx := NewContext()
	ctx.SetValue("foo", "bar")
	assert.Equal(t, ctx.GetValue("foo"), "bar")
	// Override
	ctx.SetValue("foo", 10)
	assert.Equal(t, ctx.GetValue("foo"), 10)
}

func TestContextKeyExists(t *testing.T) {
	ctx := NewContext()
	ctx.SetValue("foo", "bar")
	assert.True(t, ctx.KeyExists("foo"))
}
