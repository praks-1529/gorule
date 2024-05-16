# gorule

A powerful, lightweight and flexible rule engine written in Go. gorule allows you to define and evaluate complex business rules in an efficient and easy-to-use manner.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Rule types](#usage)
  - [Scalar rule - Scalar condition](#defining-rules)
  - [Scalar rule - Vector conditions](#evaluating-rules)
  - [Vector rules - Scalar condition](#defining-rules)
  - [Vector rules - Vector conditions](#evaluating-rules)
- [Contributing](#contributing)
- [Contact](#contact)

## Introduction

gorule is designed to simplify the process of defining and evaluating business rules within your Go applications. Whether you need to validate data, enforce policies, or build complex decision-making workflows, gorule has you covered.

## Features

- **Easy-to-Use API**: Intuitive and straightforward API for defining and managing rules.
- **High Performance**: Optimized for performance to handle large sets of rules efficiently.
- **Extensible**: Easily extendable to add custom functions and operators.
- **Comprehensive Documentation**: Detailed documentation to help you get started quickly.

## Installation

To install gorule, use `go get`:

```
go get github.com/praks-1529/gorule
```

## Quick start

```
package main

import (
    "fmt"
    "github.com/praks-1529/gorule"
)

func main() {
    // Step-1: Create a parser
    parser := gorule.NewRuleParser("IF { a == 10 && b == true && c == "Hello"} THEN { }")
    
    // Step-2: Create a rule 
	rule, err := fgParser.ParseRule()
    // Evaluate the rule

    // Step-3: Evaluate the rule across data
    testdata := []byte(`
	{
        "a" : 10,
        "b":true,
        "c":"hello",
	}
	`)

    re := NewRuleEngine()
	result, err := re.Evaluate(rule, testdata)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Rule evaluation result:", result)
}
```

## Contributing
Contributions are always welcome.

## Contact
For questions or support, please open an issue on GitHub
