// File: context.go
// Represents a Context which is used during rule evaluation
package gorule

// ContextValue represents the type of the value held inside context
type ContextValue interface{}

// Context represents the structure that is used to hold the evaluation context
type Context interface {
	SetValue(key string, value ContextValue)
	GetValue(key string) ContextValue
	KeyExists(key string) bool
}

// RuleContext represents the core context used during rule evaluation
type RuleContext struct {
	ctxMap map[string]ContextValue
}

// NewContext returns a fresh context
func NewContext() Context {
	return &RuleContext{ctxMap: make(map[string]ContextValue)}
}

// SetValue sets a value inside context
func (_ctx *RuleContext) SetValue(key string, value ContextValue) {
	_ctx.ctxMap[key] = value
}

// GetValue gets a value from the context
func (_ctx *RuleContext) GetValue(key string) ContextValue {
	return _ctx.ctxMap[key]
}

// KeyExists checks if key exists on the context
func (_ctx *RuleContext) KeyExists(key string) bool {
	_, ok := _ctx.ctxMap[key]
	return ok
}
