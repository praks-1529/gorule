//go:build !engine
// +build !engine

// File: engine_test.go
// Tests for engine
package gorule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngineScalarRuleScalarCondition(t *testing.T) {
	fgParser := NewRuleParser("IF: { domino.variantId == 3 && domino.type == \"FOO\" && domino.dpEnabled == true && domino.threshold == 5.90 }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), ScalarRuleType)
	testdata := []byte(`
	{
		"domino": {
			"variantId" : 3,
			"type":"FOO",
			"dpEnabled": true,
			"threshold": 5.90
		}
	}
	`)
	re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
	assert.Equal(t, true, result.([]bool)[0])
}

func TestEngineScalarRuleScalarConditionWithoutDomino(t *testing.T) {
	fgParser := NewRuleParser("IF: { variantId == 3 && type == \"FOO\" && dpEnabled == true && threshold == 5.90 }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), ScalarRuleType)
	testdata := []byte(`
	{
		"variantId" : 3,
		"type":"FOO",
		"dpEnabled": true,
		"threshold": 5.90
	}
	`)
	re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
	assert.Equal(t, true, result.([]bool)[0])
}

func TestEngineScalarRuleVectorConditionTrue(t *testing.T) {
	fgParser := NewRuleParser("IF: { FOR: i=0:domino.size() { domino[i].type == 10 && domino[i].dpEnabled == true } } THEN: { }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), ScalarRuleType)
	testdata := []byte(`
	{
		"domino": [{
				"type": 10,
				"dpEnabled": true,
				"threshold": 5.90
			},
			{
				"type": 10,
				"dpEnabled": true,
				"threshold": 5.90
			}
		]
	}
	`)
	re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
	assert.Equal(t, true, result.([]bool)[0])
}

func TestEngineScalarRuleVectorConditionTrueWithoutDomino(t *testing.T) {
	fgParser := NewRuleParser("IF: { FOR: i=0:b.size() { b[i].type == 10 && b[i].dpEnabled == true } } THEN: { }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), ScalarRuleType)
	testdata := []byte(`
	{
		"a" : 10,
		"b": [{
				"type": 10,
				"dpEnabled": true,
				"threshold": 5.90
			},
			{
				"type": 10,
				"dpEnabled": true,
				"threshold": 5.90
			}
		]
	}
	`)
	re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
	assert.Equal(t, true, result.([]bool)[0])
}

func TestEngineScalarRuleVectorConditionFalse(t *testing.T) {
	fgParser := NewRuleParser("IF: { FOR: i=0:domino.size() { domino[i].type == 10 } } THEN: { }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), ScalarRuleType)
	testdata := []byte(`
	{
		"domino": [{
				"type": 10,
			},
			{
				"type": 9,
			}
		]
	}
	`)
	re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
	assert.Equal(t, false, result.([]bool)[0])
}

func TestEngineVectorRuleScalarCondition(t *testing.T) {
	fgParser := NewRuleParser("FOR: i=0:domino.size() IF: { domino[i].type == 10 && domino[i].dpEnabled == true }  THEN: { }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), VectorRuleType)
	testdata := []byte(`
	{
		"domino": [{
				"type": 10,
				"dpEnabled":true
			},
			{
				"type": 9,
				"dpEnabled":true
			},
			{
				"type": 10,
				"dpEnabled":false
			},
			{
				"type": 9,
				"dpEnabled":false
			},
		]
	}
	`)
	re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
	arrResult := result.([]bool)
	assert.Equal(t, len(arrResult), 4)
	assert.Equal(t, true, arrResult[0])
	assert.Equal(t, false, arrResult[1])
	assert.Equal(t, false, arrResult[2])
	assert.Equal(t, false, arrResult[3])

}

func TestEngineVectorRuleScalarConditionWithoutDomino(t *testing.T) {
	fgParser := NewRuleParser("FOR: i=0:b.size() IF: { b[i].type == 10 && b[i].dpEnabled == true }  THEN: { }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), VectorRuleType)
	testdata := []byte(`
	{
		"a": 10,
		"b": [{
				"type": 10,
				"dpEnabled":true
			},
			{
				"type": 9,
				"dpEnabled":true
			},
			{
				"type": 10,
				"dpEnabled":false
			},
			{
				"type": 9,
				"dpEnabled":false
			},
		]
	}
	`)
	re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
	arrResult := result.([]bool)
	assert.Equal(t, len(arrResult), 4)
	assert.Equal(t, true, arrResult[0])
	assert.Equal(t, false, arrResult[1])
	assert.Equal(t, false, arrResult[2])
	assert.Equal(t, false, arrResult[3])

}
