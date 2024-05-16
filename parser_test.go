//go:build !parser
// +build !parser

// File: parser_test.go
// Tests for parser
package gorule

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperator(t *testing.T) {
	fgParser := NewRuleParser("")
	assert.True(t, fgParser.isOperator("=="))
	assert.True(t, fgParser.isOperator("&&"))
	assert.True(t, fgParser.isOperator(">="))
	assert.True(t, fgParser.isOperator("<="))
	assert.True(t, fgParser.isOperator("=="))
	assert.True(t, fgParser.isOperator(">"))
	assert.True(t, fgParser.isOperator("<"))
}

func TestPrecendence(t *testing.T) {
	fgParser := NewRuleParser("")
	assert.Equal(t, fgParser.operatorPrecendence("&&"), int16(100))
	assert.Equal(t, fgParser.operatorPrecendence("=="), int16(1))
}

func TestCheckTypeAndReturnValue(t *testing.T) {
	trueValue := StringToInterface("true")
	assert.Equal(t, trueValue, true)

	falseValue := StringToInterface("false")
	assert.Equal(t, falseValue, false)

	intValue := StringToInterface("10")
	assert.Equal(t, intValue, 10)

	floatValue := StringToInterface("10.98")
	assert.Equal(t, floatValue, 10.98)

	stringValue := StringToInterface("Test String")
	assert.Equal(t, stringValue, "Test String")
}

func TestGetNextToken(t *testing.T) {
	fgParser := NewRuleParser("valid_token")
	token, err := fgParser.getNextToken()
	assert.Nil(t, err)
	assert.Equal(t, Token("valid_token"), token)
}

func TestGetNextTokenWithHeadSpace(t *testing.T) {
	fgParser := NewRuleParser("   valid_token")
	token, err := fgParser.getNextToken()
	assert.Nil(t, err)
	assert.Equal(t, Token("valid_token"), token)
}

func TestGetNextTokenEOF(t *testing.T) {
	fgParser := NewRuleParser("valid_token")
	token, err := fgParser.getNextToken()
	assert.Nil(t, err)
	token, err = fgParser.getNextToken()
	assert.NotNil(t, err)
	assert.Equal(t, Token(""), token)
}

func TestGetNextTokenWithNewLine(t *testing.T) {
	fgParser := NewRuleParser("\nvalid_token")
	token, err := fgParser.getNextToken()
	assert.Nil(t, err)
	assert.Equal(t, Token("valid_token"), token)
}

func TestGetNextTokenWithEOF(t *testing.T) {
	fgParser := NewRuleParser("")
	_, err := fgParser.getNextToken()
	assert.NotNil(t, err)
	assert.IsType(t, &eofError{}, err)
}

func TestParserEmpty(t *testing.T) {
	fgParser := NewRuleParser("")
	_, err := fgParser.ParseRule()
	assert.NotNil(t, err)
	assert.IsType(t, &MalformedRuleError{}, err)
}

func TestParserSyntaxNoOpeningBracket(t *testing.T) {
	fgParser := NewRuleParser("IF:")
	_, err := fgParser.ParseRule()
	assert.NotNil(t, err)
	assert.IsType(t, &SyntaxError{}, err)
}

func TestParserSyntaxUnsupportedRuleType(t *testing.T) {
	fgParser := NewRuleParser("JUNK: ( )")
	_, err := fgParser.ParseRule()
	assert.NotNil(t, err)
	assert.IsType(t, &MalformedRuleError{}, err)
}

func TestParseErrorExtraOperandAtEnd(t *testing.T) {
	fgParser := NewRuleParser("IF: { a == 10 && c == true extrajunk}")
	_, err := fgParser.ParseRule()
	assert.NotNil(t, err)
	assert.IsType(t, &MalformedRuleError{}, err)
}

func TestScalarRuleScalarCondition(t *testing.T) {
	fgParser := NewRuleParser("IF: { a == 10 && c == true }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), ScalarRuleType)
	expectedRule := ScalarRule{
		Type: ScalarRuleType,
		If: &ScalarCondition{
			Type:     ScalarConditionType,
			Operator: AndOperator,
			Operand1: &ScalarCondition{
				Type:     ScalarConditionType,
				Operator: EqualOperator,
				Operand1: &ScalarCondition{
					Type:     ScalarConditionType,
					Operator: NilOperator,
					Operand1: nil,
					Operand2: nil,
					Value:    "a",
				},
				Operand2: &ScalarCondition{
					Type:     ScalarConditionType,
					Operator: NilOperator,
					Operand1: nil,
					Operand2: nil,
					Value:    10,
				},
			},
			Operand2: &ScalarCondition{
				Type:     ScalarConditionType,
				Operator: EqualOperator,
				Operand1: &ScalarCondition{
					Type:     ScalarConditionType,
					Operator: NilOperator,
					Operand1: nil,
					Operand2: nil,
					Value:    "c",
				},
				Operand2: &ScalarCondition{
					Type:     ScalarConditionType,
					Operator: NilOperator,
					Operand1: nil,
					Operand2: nil,
					Value:    true,
				},
			},
		},
	}
	actJSONData, jsonErr := json.MarshalIndent(rule, "", "\t")
	assert.Nil(t, jsonErr)
	expJSONData, jsonErr := json.MarshalIndent(expectedRule, "", "\t")
	assert.Equal(t, string(expJSONData), string(actJSONData))
}

func TestScalarRuleVectorCondition(t *testing.T) {
	fgParser := NewRuleParser("IF: { FOR: i=0:domino.size() { domino[i].type == 10 } } THEN: { }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	assert.Equal(t, rule.GetType(), ScalarRuleType)
	actJSONData, jsonErr := json.MarshalIndent(rule, "", "\t")
	assert.Nil(t, jsonErr)
	t.Log(string(actJSONData))
}

func TestVectorRuleScalarCondition(t *testing.T) {
	fgParser := NewRuleParser("FOR: i=0:domino.size() IF: { domino[i].type == 1 && domino[i+1] == 0 }  THEN: { }")
	rule, err := fgParser.ParseRule()
	assert.Nil(t, err)
	actJSONData, jsonErr := json.MarshalIndent(rule, "", "\t")
	assert.Nil(t, jsonErr)
	t.Log(string(actJSONData))
}
