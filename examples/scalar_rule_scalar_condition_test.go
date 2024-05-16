package examples

import (
	"fmt"
	"testing"

	"github.com/praks-1529/gorule"
)

// Given a transaction data, categorizes it as risky or not based on the amount and type of transaction
func TestScalarRuleScalarConditionExample(t *testing.T) {
	// Step-1: Create a parser
	parser := gorule.NewRuleParser("IF: { amount >= 10000 && type == \"CREDIT_CARD\" }")

	// Step-2: Create a rule
	rule, err := parser.ParseRule()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Step-3: Evaluate the rule across data
	riskyTxn := []byte(`
	{
        "amount" : 10000,
		"type":"CREDIT_CARD",
	}
	`)
	nonRiskyTxn := []byte(`
	{
        "amount" : 9999,
		"type":"CREDIT_CARD",
	}
	`)

	re := gorule.NewRuleEngine()
	result1, err := re.Evaluate(rule, riskyTxn)
	result2, err := re.Evaluate(rule, nonRiskyTxn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Rule evaluation result:", result1, result2)
}
