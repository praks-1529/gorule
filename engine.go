// File: engine.go
// Rule engine to evaluate the rule for the given data
package gorule

// RuleEngine represents a service to evaluate the rule
type RuleEngine struct {
}

// NewRuleEngine returns a fresh rule engine instance
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (_re *RuleEngine) buildContext(fgRule Rule, ipData []byte) Context {
	ctx := NewContext()
	fgRule.BuildContext(ipData, ctx)
	return ctx
}

// Evaluate evaluates a rule for the given rule and jsonData
// args:
//
//	fgRule: The rule to evaluate
//	jsonData: The data to be used during evaluation
//
// Return
//
//	bool: Evaluation result (i.e true/false)
//	error: Any error during evaluation
func (_re *RuleEngine) Evaluate(fgRule Rule, jsonData []byte) (interface{}, error) {
	ctx := _re.buildContext(fgRule, jsonData)
	return fgRule.Evaluate(ctx)
}
