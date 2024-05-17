package examples

import (
	"fmt"
	"testing"

	"github.com/praks-1529/gorule"
)

// Below example iterates over each transactions and categorizes the transactions as risky or not based on the amount and type of transaction
// This is same as running scalar-rule X N times where N is the number of transactions
func TestVectorRuleScalarConditionExample(t *testing.T) {
	parser := gorule.NewRuleParser("FOR: i=0:transactions.size() IF: { transactions[i].amount > 10000 && transactions[i].type == \"CREDIT_CARD\" }")

	// Step-2: Create a rule
	rule, err := parser.ParseRule()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Step-3: Evaluate the rule across data
	transactions := []byte(`
	{
		"transactions": [
		  {
			"id": 1,
			"some_other_payment_meta": {},
			"type": "CREDIT_CARD",
			"amount": 10000
		  },
		  {
			"id": 2,
			"some_other_payment_meta": {},
			"type": "CREDIT_CARD",
			"amount": 10001
		  },
		  {
			"id": 3,
			"some_other_payment_meta": {},
			"type": "DEBIT_CARD",
			"amount": 10001
		  }
		]
	}
	`)
	re := gorule.NewRuleEngine()
	isRiskyTransactions, _ := re.Evaluate(rule, transactions)
	fmt.Println("Rule evaluation result:", isRiskyTransactions)
}
