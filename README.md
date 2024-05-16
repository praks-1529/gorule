# gorule

A powerful, lightweight and flexible rule engine written in Go. gorule allows you to define and evaluate complex business rules in an efficient and easy-to-use manner.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Rule types](#rule-types)
  - [Scalar rule - Scalar condition](#scalar-rule-scalar-condition)
  - [Scalar rule - Vector conditions](#scalar-rule-vector-condition)
  - [Vector rules - Scalar condition](#vector-rule-vector-condition)
- [Supported data types](#supported-data-types)
- [Supported operators](#supported-operators)
- [Contributing](#contributing)
- [Contact](#contact)

## Introduction

gorule is designed to simplify the process of defining and evaluating business rules within your Go applications. Whether you need to validate data, enforce policies, or build complex decision-making workflows, gorule has you covered.

## Features

- **Easy-to-Use**: Define the rules in human readable form
- **High Performance**: Optimized for performance to handle large sets of rules efficiently. 
- **Extensible**: Rule engine that works for deeply nested JSON objects INCLUDING array fields (i.e vector conditions support)
- **Comprehensive Documentation**: Detailed documentation to help you get started quickly.

## Installation

To install gorule, use `go get`:

```
go get github.com/praks-1529/gorule
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/praks-1529/gorule"
)

// Given a transaction data, categorizes it as risky or not based on the amount and type of transaction
func main() {
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
```
Output

```sh
Rule evaluation result: [true] [false]
```

See [examples] for more advanced use cases

## Rule Types

`gorule` supports various types of rules to cater to different use cases. Below are the rule types supported by `gorule`:

### Scalar Rule Scalar Condition
Scalar rules with scalar condition evaluate a single set of condition on a JSON object. These kind of rules are applicable for simple JSON objects


#### Example
```go
// Define a rule to categorize a transaction as risky if the amount is greater than or equal to 10000 and the type is "CREDIT_CARD"
parser := gorule.NewRuleParser("IF: { transaction.amount >= 10000 && transaction.type == \"CREDIT_CARD\" }")
```

#### Example use case
Given a transaction data, categorize it as risky or not based on the amount and type of transaction

### Scalar Rule Vector Condition
Scalar rules with vector conditions evaluate conditions over an array of attributes within a single JSON object. These rules are useful when the decision depends on array fields inside JSON object.

Example

```go
// Below example categorizes the transactions as risky or not based on multiple attributes of the transaction
parser := gorule.NewRuleParser("IF: { FOR: i=0:transaction.attributes.size() { transaction.attributes[i].type == \"VALID\" } }")

```

#### Example use case
Categorize transactions as risky or not based on the attributes of the transaction.


### Vector Rule Vector Condition
These rules are useful when the decision depends on iterating over and evaluating a set of JSON objects in addition to decision depending on array fields inside each JSON object.

Example

```go
// Below example iterates over each transactions and categorizes the transactions as risky or not based on the amount and type of transaction
parser := gorule.NewRuleParser("FOR: i=0:transactions.size() IF: { transactions[i].amount > 10000 && transactions[i].type == \"CREDIT_CARD\" }")
```

#### Example use case
Iterate over each transaction and categorize the transactions as risky or not based on the amount and type of transaction. 

## Supported data types
- Integers
- Float
- String
- Boolean

## Supported operators
| Operator | Description | Precendence |
| -------- | ------- | -------| 
| && | Logical AND | 2 |
| \|\| | Logical OR | 2 |
| == | Equal to | 1 |
| >= | Greater than or equal to | 1 |
| > | Greater than | 1 |
| >= | Lesser than or equal to | 1 |
| > | Lesser than | 1 |

## Contributing
Contributions are always welcome.

## Contact
For questions or support, please open an issue on GitHub
