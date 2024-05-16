// File: constants.go
// Stores all the constant shared within packages
package gorule

// ConditionType represents the condition type
type ConditionType int

const (
	// IndexKey ...
	IndexKey string = "_INDEX_KEY"
	// IndexCurrentValue ..
	IndexCurrentValue string = "_INDEX_CURRENT_VALUE"
	// StartIndexValue ...
	StartIndexValue string = "_START_INDEX_VALUE"
	// EndIndexValue ...
	EndIndexValue string = "_END_INDEX_VALUE"
)
