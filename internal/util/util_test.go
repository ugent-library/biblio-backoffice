package util

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

// Test ParseString
func TestParseString(t *testing.T) {
	assert.Equal(t, ParseString("test"), "test")
	assert.Equal(t, ParseString(true), "true")
	assert.Equal(t, ParseString(int(1)), "1")
	assert.Equal(t, ParseString(int32(2_000_000_000)), "2000000000")
	assert.Equal(t, ParseString(int64(8_000_000_000)), "8000000000")
	assert.Equal(t, ParseString(float32(1.2)), "1.2")
	assert.NotEqual(t, ParseString(float32(1.20)), "1.20")
	assert.Equal(t, ParseString(float64(1.2)), "1.2")
	var i interface{}
	assert.Equal(t, ParseString(i), "")
	assert.Equal(t, ParseString(nil), "")
}

func TestParseBoolean(t *testing.T) {
	assert.Equal(t, ParseBoolean(true), true)
	assert.Equal(t, ParseBoolean(false), false)
	assert.Equal(t, ParseBoolean("true"), true)
	assert.Equal(t, ParseBoolean("false"), false)
	assert.Equal(t, ParseBoolean(int32(1)), true)
	assert.Equal(t, ParseBoolean(int32(0)), false)
	assert.Equal(t, ParseBoolean(int64(1)), true)
	assert.Equal(t, ParseBoolean(int64(0)), false)
}
