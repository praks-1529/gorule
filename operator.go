// File: operator.go
// Defines all the supported operators
package gorule

import (
	"reflect"
)

// Operator represents the type of operator
type Operator string

const (
	// AndOperator for representing logical AND (works for bool ONLY)
	AndOperator Operator = "&&"
	// OrOperator for representing logical OR (works on bool ONLY)
	OrOperator Operator = "||"
	// EqualOperator for representing logical == (works for int, string, bool. float)
	EqualOperator Operator = "=="
	// GreaterThanOrEqualOperator for representing logical >= (works for int, float)
	GreaterThanOrEqualOperator Operator = ">="
	// GreaterOperator for representing logical >= (works for int, float)
	GreaterOperator Operator = ">"
	// LesserThanOrEqualOperator for representing logical <= (works for int, float)
	LesserThanOrEqualOperator Operator = "<="
	// LesserOperator for representing logical <= (works for int, float)
	LesserOperator Operator = "<"
	// NilOperator for representing NIL operator
	NilOperator Operator = "NIL"
)

func evaluateInt(op1 int, op2 int, optor Operator) bool {
	switch optor {
	case EqualOperator:
		return op1 == op2
	case GreaterOperator:
		return op1 > op2
	case GreaterThanOrEqualOperator:
		return op1 >= op2
	case LesserOperator:
		return op1 < op2
	case LesserThanOrEqualOperator:
		return op1 <= op2
	}
	return false
}

func evaluateFloat64(op1 float64, op2 float64, optor Operator) bool {
	switch optor {
	case EqualOperator:
		return op1 == op2
	case GreaterOperator:
		return op1 > op2
	case GreaterThanOrEqualOperator:
		return op1 >= op2
	case LesserOperator:
		return op1 < op2
	case LesserThanOrEqualOperator:
		return op1 <= op2
	}
	return false
}

func evaluateString(op1 string, op2 string, optor Operator) bool {
	switch optor {
	case EqualOperator:
		return op1 == op2
	}
	return false
}

func evaluateBool(op1 bool, op2 bool, optor Operator) bool {
	switch optor {
	case AndOperator:
		return op1 && op2
	case OrOperator:
		return op1 || op2
	case EqualOperator:
		return op1 == op2
	}
	return false
}

// EvaluateOperation is used to evaluate supported operations
// TODO: To make this more generic based on reflect package
func EvaluateOperation(operand1 interface{}, operand2 interface{}, optor Operator) bool {
	if reflect.TypeOf(operand1) != reflect.TypeOf(operand2) {
		panic("Operands type not matching")
	}
	switch op1 := operand1.(type) {
	case int:
		return evaluateInt(op1, operand2.(int), optor)
	case float64:
		return evaluateFloat64(op1, operand2.(float64), optor)
	case string:
		return evaluateString(op1, operand2.(string), optor)
	case bool:
		return evaluateBool(op1, operand2.(bool), optor)
	default:
		panic("Unsupported type found")
	}
}
