//go:build !parser
// +build !parser

// File: engine_test.go
// Tests for engine
package gorule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToInterface(t *testing.T) {
	assert.Equal(t, true, StringToInterface("true"))
	assert.Equal(t, false, StringToInterface("false"))
	assert.Equal(t, int(1), StringToInterface("1"))
	assert.Equal(t, float64(1.2), StringToInterface("1.2"))
}

func TestResolveValue(t *testing.T) {
	testData := []byte(`
	{
		"domino": {
			"variantId" : 3,
			"type":"FOO",
			"dpEnabled": true,
			"pdpDisabled": false,
			"threshold": 5.90,
			"moves": [
				{
					"type":100,
					"user":1
				},
				{
					"type":1,
					"user":2
				}
			]
		}
	}
	`)
	// Test integer
	assert.Equal(t, int(3), resolveValue("domino.variantId", testData))
	// Test string
	assert.Equal(t, "\"FOO\"", resolveValue("domino.type", testData))
	// Test float
	assert.Equal(t, float64(5.90), resolveValue("domino.threshold", testData))
	// Test float-1
	assert.Equal(t, true, resolveValue("domino.dpEnabled", testData))
	// Test float-2
	assert.Equal(t, false, resolveValue("domino.pdpDisabled", testData))
	// Test fetching from array index
	assert.Equal(t, int(100), resolveValue("domino.moves.0.type", testData))
}

func TestResolveLength(t *testing.T) {
	testData := []byte(`
	{
		"domino": {
			"moves": [
				{
					"type":1,
					"user":1,
					"dummy":[1,2,3,4,5,6,7]
				},
				{
					"type":1,
					"user":2,
					"dummy":[1,2,3,4,5,6,7]
				}
			]
		}
	}
	`)
	// Test integer
	assert.Equal(t, int(2), resolveLength("domino.moves", testData))
}

func TestHasArrayIndex(t *testing.T) {
	assert.Equal(t, false, hasArrayIndex("abc"))
	assert.Equal(t, true, hasArrayIndex("abc[i]"))

}

func TestResolveValue2(t *testing.T) {
	testData := []byte(`
	{
		"variantId" : 3,
		"type":"FOO",
		"dpEnabled": true,
		"pdpDisabled": false,
		"threshold": 5.90,
		"moves": [
			{
				"type":100,
				"user":1
			},
			{
				"type":1,
				"user":2
			}
		]
	}
	`)
	// Test integer
	assert.Equal(t, int(3), resolveValue("variantId", testData))
	// Test string
	assert.Equal(t, "\"FOO\"", resolveValue("type", testData))
	// Test float
	assert.Equal(t, float64(5.90), resolveValue("threshold", testData))
	// Test float-1
	assert.Equal(t, true, resolveValue("dpEnabled", testData))
	// Test float-2
	assert.Equal(t, false, resolveValue("pdpDisabled", testData))
	// Test fetching from array index
	assert.Equal(t, int(100), resolveValue("moves.0.type", testData))
}
