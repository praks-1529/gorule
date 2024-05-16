// File: condition.go
// Represents one condition ex. a == b
package gorule

import (
	"fmt"
	"strings"
)

const (
	// ScalarConditionType represents a binary condition which is one dimensional (ex. a == b && c == d)
	ScalarConditionType ConditionType = 1
	// VectorConditionType represents iterative condition (ex. for i=0:N (result = result && a[i]))
	VectorConditionType ConditionType = 2
)

// Condition represents the interface for condition
type Condition interface {
	GetOperator() Operator
	Evaluate(ctx Context) (interface{}, error)
	GetValue() interface{}
	buildContext(ipData []byte, ctx Context)
}

// ScalarCondition Condition represents one evaluatable binary expression (ex: a == b)
// Format: { a == b && c == d ) }
type ScalarCondition struct {
	Type          ConditionType `json:"type"`
	Operator      Operator      `json:"optor"`
	Operand1      Condition     `json:"o1"`
	Operand2      Condition     `json:"o2"`
	Value         interface{}   `json:"value"`
	HasArrayIndex bool          `json:"has_array_index"`
}

// Evaluate does the evaluation of the condition and returns the result
func (_c *ScalarCondition) Evaluate(ctx Context) (interface{}, error) {
	if _c.GetOperator() == NilOperator {
		if _, ok := _c.GetValue().(string); ok {
			// If its is a string literal it can be a json field like "a.b"
			// or it can be a pure string like "FOO". In both the cases we add them
			// in the context
			ctxKey := _c.getContextKey(ctx)
			return ctx.GetValue(ctxKey), nil
		}
		return _c.GetValue(), nil
	}
	lvalue, _ := _c.GetOperand1().Evaluate(ctx)
	rvalue, _ := _c.GetOperand2().Evaluate(ctx)
	result := EvaluateOperation(lvalue, rvalue, _c.GetOperator())
	return result, nil
}

func (_c *ScalarCondition) getArrayIndex(ctx Context) string {
	key := fmt.Sprintf("[%s]", ctx.GetValue(IndexKey))
	return key
}

func (_c *ScalarCondition) getParsableArrayIndex(ctx Context) string {
	key := fmt.Sprintf(".%d", ctx.GetValue(IndexCurrentValue))
	return key
}

func (_c *ScalarCondition) getContextKey(ctx Context) string {
	key := _c.GetValue().(string)
	if _c.HasArrayIndex {
		key = strings.Replace(key, _c.getArrayIndex(ctx), _c.getParsableArrayIndex(ctx), 1)
	}
	return key
}

func (_c *ScalarCondition) buildContext(ipData []byte, ctx Context) {
	if _c.GetOperator() == NilOperator {
		if _, ok := _c.GetValue().(string); ok {
			ctxKey := _c.getContextKey(ctx)
			value := resolveValue(ctxKey, ipData)
			ctx.SetValue(ctxKey, value)
		}
		return
	}
	_c.GetOperand1().buildContext(ipData, ctx)
	_c.GetOperand2().buildContext(ipData, ctx)
}

// GetOperator returns the underlying operator in the condition
func (_c *ScalarCondition) GetOperator() Operator {
	return _c.Operator
}

// GetOperand1 returns the underlying operand-1
func (_c *ScalarCondition) GetOperand1() Condition {
	return _c.Operand1
}

// GetOperand2 returns the underlying operand-2
func (_c *ScalarCondition) GetOperand2() Condition {
	return _c.Operand2
}

// GetValue returns data stored in the condition
func (_c *ScalarCondition) GetValue() interface{} {
	return _c.Value
}

// VectorCondition represents collection of Conditions joined by operator
// Format: { FOR i=initialValue:finalValue SCALAR_CONDITION }
type VectorCondition struct {
	Type       ConditionType `json:"type"`
	Operator   Operator      `json:"optor"`
	SCondition Condition     `json:"scalar_condition"`
	Value      interface{}   `json:"value"`
	StartIndex interface{}   `json:"start_index"`
	EndIndex   interface{}   `json:"end_index"`
	IndexKey   string        `json:"index_key"`
}

// GetOperator returns the underlying operator in the condition
func (_c *VectorCondition) GetOperator() Operator {
	return _c.Operator
}

// Evaluate does the evaluation of the condition and returns the result
func (_c *VectorCondition) Evaluate(ctx Context) (interface{}, error) {
	result := true
	ctx.SetValue(IndexKey, _c.IndexKey)
	for i := ctx.GetValue(StartIndexValue).(int); i < ctx.GetValue(EndIndexValue).(int); i++ {
		var res interface{}
		var err error
		scalarCondition := _c.SCondition.(*ScalarCondition)
		ctx.SetValue(IndexCurrentValue, i)
		// TODO: The default operator is &&
		if res, err = scalarCondition.Evaluate(ctx); err != nil {
			return false, err
		}
		// TODO: Move this to Result structure
		result = result && res.(bool)
	}
	_c.Value = result
	return _c.GetValue(), nil
}

// GetValue returns data stored in the condition
func (_c *VectorCondition) GetValue() interface{} {
	return _c.Value
}

// Evaluate string like "0"
func (_c *VectorCondition) getInitialValue(ipData []byte) int {
	value := StringToInterface(_c.StartIndex.(string)).(int)
	return value
}

// Evaluate string like a.size() => len(a)
func (_c *VectorCondition) getFinalValue(ipData []byte) int {
	value := resolveLength(strings.Replace(_c.EndIndex.(string), ".size()", "", 1), ipData).(int)
	return value
}

func (_c *VectorCondition) buildContext(ipData []byte, ctx Context) {
	startIndex := _c.getInitialValue(ipData)
	ctx.SetValue(StartIndexValue, startIndex)
	endIndex := _c.getFinalValue(ipData)
	ctx.SetValue(EndIndexValue, endIndex)
	// Set this IndexKey = "i"
	ctx.SetValue(IndexKey, _c.IndexKey)
	// For every value of "i", fetch the value
	for i := startIndex; i < endIndex; i++ {
		// Set this "[i] = 0"
		ctx.SetValue(IndexCurrentValue, i)
		scalarCondition := _c.SCondition.(*ScalarCondition)
		scalarCondition.buildContext(ipData, ctx)
	}

}
