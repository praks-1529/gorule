// File: rule.go
// Implement the rule interface
package gorule

import "strings"

// RuleType represents types of rule
type RuleType int

const (
	// ScalarRuleType represents a simple rule ex. IF: (a==b && c==d)) THEN: ()
	ScalarRuleType RuleType = 1
	// VectorRuleType represents iterative rule ex: FOR i=0:a.size(): IF (a[i].type == a) THEN: ()
	VectorRuleType RuleType = 2
)

// Rule represent the common interface for any rule
type Rule interface {
	GetType() RuleType
	Evaluate(ctx Context) (interface{}, error)
	BuildContext(ipData []byte, ctx Context) error
}

// ScalarRule represents a structure of the If rule (created by parsing the user provided rules)
type ScalarRule struct {
	Type       RuleType    `json:"type"`
	If         Condition   `json:"if"`
	Then       Action      `json:"then"`
	StartIndex interface{} `json:"start_index"`
	IndexKey   interface{} `json:"index_key"`
}

// Evaluate evalates the rule
func (_fgr *ScalarRule) Evaluate(ctx Context) (interface{}, error) {
	var result []bool
	res, err := _fgr.If.Evaluate(ctx)
	result = append(result, res.(bool))
	return result, err
}

// GetType gets the type of rule
func (_fgr *ScalarRule) GetType() RuleType {
	return _fgr.Type
}

func (_fgr *ScalarRule) getCondition() Condition {
	return _fgr.If
}

// BuildContext builds the context
func (_fgr *ScalarRule) BuildContext(ipData []byte, ctx Context) error {
	_fgr.getCondition().buildContext(ipData, ctx)
	// TODO: Above condition can return errors
	return nil
}

// VectorRule represents a collection of Simple rules
type VectorRule struct {
	Type       RuleType    `json:"type"`
	SRule      Rule        `json:"scalar_rule"`
	StartIndex interface{} `json:"start_index"`
	EndIndex   interface{} `json:"end_index"`
	IndexKey   interface{} `json:"index_key"`
}

// Evaluate evalates the rule
func (_fgr *VectorRule) Evaluate(ctx Context) (interface{}, error) {
	var result []bool
	ctx.SetValue(IndexKey, _fgr.IndexKey)
	for i := ctx.GetValue(StartIndexValue).(int); i < ctx.GetValue(EndIndexValue).(int); i++ {
		var res interface{}
		var err error
		scalarRule := _fgr.SRule.(*ScalarRule)
		ctx.SetValue(IndexCurrentValue, i)
		// TODO: The default operator is &&
		if res, err = scalarRule.Evaluate(ctx); err != nil {
			result = append(result, false)
		}
		result = append(result, res.([]bool)[0])
	}
	return result, nil
}

// GetType gets the type of rule
func (_fgr *VectorRule) GetType() RuleType {
	return _fgr.Type
}

// Evaluate string like "0"
func (_fgr *VectorRule) getInitialValue(ipData []byte) int {
	value := StringToInterface(_fgr.StartIndex.(string)).(int)
	return value
}

// Evaluate string like a.size() => len(a)
func (_fgr *VectorRule) getFinalValue(ipData []byte) int {
	value := resolveLength(strings.Replace(_fgr.EndIndex.(string), ".size()", "", 1), ipData).(int)
	return value
}

// BuildContext builds the context
func (_fgr *VectorRule) BuildContext(ipData []byte, ctx Context) error {
	startIndex := _fgr.getInitialValue(ipData)
	ctx.SetValue(StartIndexValue, startIndex)
	endIndex := _fgr.getFinalValue(ipData)
	ctx.SetValue(EndIndexValue, endIndex)
	// Set this IndexKey = "i"
	ctx.SetValue(IndexKey, _fgr.IndexKey)
	// For every value of "i", fetch the value
	for i := startIndex; i < endIndex; i++ {
		// Set this "[i] = 0"
		ctx.SetValue(IndexCurrentValue, i)
		scalarRule := _fgr.SRule.(*ScalarRule)
		scalarRule.BuildContext(ipData, ctx)
	}
	return nil
}
