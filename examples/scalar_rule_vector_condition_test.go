package examples

import (
	"fmt"
	"testing"

	"github.com/praks-1529/gorule"
)

// Below example categorizes the transactions as risky or not based on attributes of the transaction
// This is same as running scalar-rule with N X condition times where N is the number of conditions
func TestScalarRuleVectorConditionExample(t *testing.T) {
	parser := gorule.NewRuleParser("IF: { FOR: i=0:transaction.attributes.size() { transaction.attributes[i].type == \"VALID\" } }")

	// Step-2: Create a rule
	rule, err := parser.ParseRule()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Step-3: Evaluate the rule across data
	riskyTxn := []byte(`
	{
		"transaction": {
		  "amount": 10000,
		  "type": "CREDIT_CARD",
		  "attributes": [
			{
			  "key": "IP",
			  "value": "10.00.00.00",
			  "type": "FRAUD"
			},
			{
			  "key": "DEVICE",
			  "value": "iPhone",
			  "type": "VALID"
			},
			{
			  "key": "LOCATION",
			  "value": "ZZ",
			  "type": "VALID"
			}
		  ]
		}
	}
	`)
	nonRiskyTxn := []byte(`
	{
		"amount": 10000,
		"type": "CREDIT_CARD",
		"transaction": {
		  "attributes": [
			{
			  "key": "IP",
			  "value": "10.11.11.12",
			  "type": "VALID"
			},
			{
			  "key": "DEVICE",
			  "value": "iPhone",
			  "type": "VALID"
			},
			{
			  "key": "LOCATION",
			  "value": "ZZ",
			  "type": "VALID"
			}
		  ]
		}
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
