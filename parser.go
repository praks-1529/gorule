// File: parser.go
// Implements the rule parser
package gorule

import (
	"fmt"
	"strings"

	"github.com/golang-collections/collections/stack"
)

// Token represents a valid token.
type Token string

const (
	// IfToken represents start of If rule
	IfToken Token = "IF:"
	// ThenToken represents start of Then block
	ThenToken Token = "THEN:"
	// OpenBraceToken represents start of block
	OpenBraceToken Token = "("
	// CloseBraceToken represents close of block
	CloseBraceToken Token = ")"
	// CurlyOpenBraceToken represents start of a condition/action
	CurlyOpenBraceToken Token = "{"
	// CurlyCloseBraceToken represents close of block
	CurlyCloseBraceToken Token = "}"
	// ColonToken represents delimiter
	ColonToken Token = ":"
	// ForToken represents start if For rule
	ForToken Token = "FOR:"
)

// EOFError rerpresents end of file error
type eofError struct {
}

func (_rt *eofError) Error() string {
	return "End of file reached"
}

// SyntaxError raised when there is a syntax error
type SyntaxError struct {
	Expected Token
	Found    Token
	Index    int
}

func (_rt *SyntaxError) Error() string {
	return fmt.Sprintf("Syntax error : Expected %s found %s", _rt.Expected, _rt.Found)
}

// MalformedRuleError raised when the input data is malformed in some way
type MalformedRuleError struct {
}

func (_rt *MalformedRuleError) Error() string {
	return fmt.Sprintf("Rule format is malformed")
}

// RuleParser service to parse the rule
type RuleParser struct {
	currentToken Token
	prevToken    Token
	currentIndex int
	optorStack   *stack.Stack
	oprndStack   *stack.Stack
	input        string
}

// NewRuleParser returns the fresh instance of RuleParser
func NewRuleParser(ip string) *RuleParser {
	parserInstance := &RuleParser{}
	parserInstance.init(ip)
	return parserInstance
}

func (_p *RuleParser) init(ip string) {
	_p.optorStack = stack.New()
	_p.oprndStack = stack.New()
	_p.currentIndex = 0
	_p.input = ip
}

func (_p *RuleParser) isTokenValid(token Token) bool {
	return false
}

// Returns the next token of the parser
// TODO: Currently works on space delimiter, but must be extended to support based on tokens
func (_p *RuleParser) getNextToken() (Token, error) {
	if _p.currentIndex >= len(_p.input) {
		return "", &eofError{}
	}
	var token Token
	token = ""
	// Trim all white spaces/new lines at head
	for _p.currentIndex < len(_p.input) && (_p.input[_p.currentIndex] == ' ' || _p.input[_p.currentIndex] == '\n') {
		_p.currentIndex++
	}
	// Check if we have reached end of file
	if _p.currentIndex >= len(_p.input) {
		// End of file reached
		return "", &eofError{}
	}
	// Exit if we run out of string OR we find space OR token so far is valid
	for !(_p.currentIndex >= len(_p.input) || _p.input[_p.currentIndex] == ' ' || _p.isTokenValid(token)) {
		token += Token(_p.input[_p.currentIndex])
		_p.currentIndex++
	}
	_p.prevToken = _p.currentToken
	_p.currentToken = token
	return _p.currentToken, nil
}

func (_p *RuleParser) getCurrentIndex() int {
	return _p.currentIndex
}

func (_p *RuleParser) rewind(rewindIndex int) {
	_p.currentIndex = rewindIndex
}

// Returns the previous token of the parser
func (_p *RuleParser) getPrevToken() Token {
	return _p.prevToken
}

func (_p *RuleParser) isOperator(token Token) bool {
	if len(token) <= 2 && (token == Token(AndOperator) ||
		token == Token(OrOperator) ||
		token == Token(EqualOperator) ||
		token == Token(GreaterOperator) ||
		token == Token(GreaterThanOrEqualOperator) ||
		token == Token(LesserOperator) ||
		token == Token(LesserThanOrEqualOperator)) {
		return true
	}
	return false
}

func (_p *RuleParser) operatorPrecendence(token Operator) int16 {
	if !_p.isOperator(Token(token)) {
		panic(token)
	}
	switch token {
	case AndOperator:
		return 100
	case OrOperator:
		return 100
	default:
		return 1
	}
}

func (_p *RuleParser) formExpression() Condition {
	// Take top two items frm opnd stack and mix with topOptor
	stkTop := _p.optorStack.Peek()
	if topOptor, ok := stkTop.(Operator); ok {
		op2, _ := _p.oprndStack.Pop().(Condition)
		op1, _ := _p.oprndStack.Pop().(Condition)
		curCond := &ScalarCondition{Type: ScalarConditionType, Operator: topOptor, Value: nil, Operand1: op1, Operand2: op2}
		_p.oprndStack.Push(curCond)
		_p.optorStack.Pop()
		return curCond
	}
	return nil
}

func (_p *RuleParser) createLeafCond(token interface{}) Condition {
	hasArrIndex := false
	_, ok := token.(string)
	if ok && hasArrayIndex(token.(string)) {
		hasArrIndex = true
	}
	return &ScalarCondition{Type: ScalarConditionType, Operator: NilOperator, Value: token, Operand1: nil, Operand2: nil, HasArrayIndex: hasArrIndex}
}

// Parse the token such as i=0:b.size() and return index key, start_index(string), end_index(string)
func (_p *RuleParser) parseForVectorDefinitions() (string, interface{}, interface{}, error) {
	token, err := _p.getNextToken()
	if err != nil {
		return "nil", nil, nil, err
	}
	subTokens := strings.Split(string(token), "=")
	indexKey := subTokens[0]
	subTokens = strings.Split(subTokens[1], ":")
	return indexKey, subTokens[0], subTokens[1], nil
}

func (_p *RuleParser) validateConditionStart() error {
	curToken, err := _p.getNextToken()
	// Conditions always start with {
	if curToken != CurlyOpenBraceToken || err != nil {
		return &SyntaxError{Expected: CurlyOpenBraceToken, Found: curToken, Index: _p.currentIndex}
	}
	return nil
}

// Parse scalar condition and create a conditions tree
//
// Format: { a == b && c == d ) }
//
// Returns the parsed condition
func (_p *RuleParser) parseScalarCondition() (Condition, error) {
	// Conditions always start with {
	if err := _p.validateConditionStart(); err != nil {
		return nil, err
	}
	var curToken Token
	var err error
	curToken, err = _p.getNextToken()
	for curToken != CurlyCloseBraceToken && err == nil {
		if !_p.isOperator(curToken) {
			// Leaf level node in the decision tree
			parsedValue := StringToInterface(string(curToken))
			leafCond := _p.createLeafCond(parsedValue)
			_p.oprndStack.Push(leafCond)
		} else {
			curOptor := Operator(curToken)
			if _p.optorStack.Len() == 0 {
				// Push the opperator
				_p.optorStack.Push(curOptor)
				curToken, err = _p.getNextToken()
				continue
			}
			stkTop := _p.optorStack.Peek()
			// fmt.Println("Parsed token : " + tokens[i])
			if topOptor, ok := stkTop.(Operator); ok {
				if _p.operatorPrecendence(curOptor) > _p.operatorPrecendence(topOptor) {
					// Take top two items frm opnd stack and mix with topOptor
					_ = _p.formExpression()
					_p.optorStack.Push(curOptor)
				} else {
					// Push the opperator
					_p.optorStack.Push(curOptor)
				}
			}
		}
		curToken, err = _p.getNextToken()
	}
	for _p.optorStack.Len() > 0 {
		_ = _p.formExpression()
	}
	if _p.oprndStack.Len() > 1 {
		return nil, &MalformedRuleError{}
	}
	return _p.oprndStack.Pop().(Condition), nil
}

// Parse vector condition and create a conditions tree
//
// Format: { FOR i=initialValue:finalValue SCALAR_CONDITION }
//
// Returns the parsed condition
func (_p *RuleParser) parseVectorCondition() (Condition, error) {
	// Conditions always start with {
	if err := _p.validateConditionStart(); err != nil {
		return nil, err
	}
	var curToken Token
	var err error
	// TODO: The operator is hard coded to &&. Can be extended to anything
	vectorCondition := &VectorCondition{Type: VectorConditionType, Operator: AndOperator}
	if curToken, err = _p.getNextToken(); curToken != ForToken || err != nil {
		return nil, &SyntaxError{Expected: ForToken, Found: curToken, Index: _p.currentIndex}
	}
	if vectorCondition.IndexKey, vectorCondition.StartIndex, vectorCondition.EndIndex, err = _p.parseForVectorDefinitions(); err != nil {
		return nil, err
	}
	if vectorCondition.SCondition, err = _p.parseCondition(); err != nil {
		return nil, err
	}
	return vectorCondition, nil
}

func (_p *RuleParser) parseCondition() (Condition, error) {
	rewindIndex := _p.currentIndex
	var err error
	if err = _p.validateConditionStart(); err != nil {
		// Conditions always start with {
		return nil, err
	}
	var curToken Token
	if curToken, err = _p.getNextToken(); err == nil {
		// Rewind so that it starts either at FOR or rule start
		_p.rewind(rewindIndex)
		if curToken == ForToken {
			return _p.parseVectorCondition()
		}
		return _p.parseScalarCondition()
	}
	return nil, err
}

// Validates that rules start with either IF: or FOR:
func (_p *RuleParser) validateRuleStart() error {
	curToken, err := _p.getNextToken()
	// Conditions always start with {
	if !(curToken == IfToken || curToken == ForToken) || err != nil {
		return &SyntaxError{Expected: IfToken, Found: curToken, Index: _p.currentIndex}
	}
	return nil
}

// Implement parser the scalar rule
//
// Format: IF: { CONDITION } THEN: { ACTION }
//
// Returns the parsed condition
func (_p *RuleParser) parseScalarRule() (Rule, error) {
	if err := _p.validateRuleStart(); err != nil {
		return nil, err
	}
	retRule := &ScalarRule{Type: ScalarRuleType}
	var err error
	// Parse the condition
	retRule.If, err = _p.parseCondition()
	if err != nil {
		return nil, err
	}
	// Parse the action
	retRule.Then, err = _p.parseActions()
	if err != nil {
		return nil, err
	}
	return retRule, err
}

// Implement parser to parse vector rule
//
// Format: FOR: index=initialValue:finalValue SCALAR_RULE
//
//	FOR i=initialValue:finalValue: SCALAR_RULE
func (_p *RuleParser) parseVectorRule() (Rule, error) {
	var err error
	if err = _p.validateRuleStart(); err != nil {
		return nil, err
	}
	vectorRule := &VectorRule{Type: VectorRuleType}
	if vectorRule.IndexKey, vectorRule.StartIndex, vectorRule.EndIndex, err = _p.parseForVectorDefinitions(); err != nil {
		return nil, err
	}
	if vectorRule.SRule, err = _p.parseScalarRule(); err != nil {
		return nil, err
	}
	return vectorRule, nil
}

func (_p *RuleParser) parseActions() (Rule, error) {
	return nil, nil
}

// ParseRule main entry to parse the rule
// args:
//
//	s (string): The input rule to parse
//
// returns
// Rule: Parsed rule
// error: Error any dound while parsing
// TODO: Handle all cases - For now simplest case a && b && c
func (_p *RuleParser) ParseRule() (Rule, error) {
	rewindIndex := _p.currentIndex
	// Outer layer has to be one of IF: / FOR:
	ruleToken, e := _p.getNextToken()
	if e != nil {
		return nil, &MalformedRuleError{}
	}
	_p.rewind(rewindIndex)
	switch ruleToken {
	case IfToken:
		return _p.parseScalarRule()
	case ForToken:
		return _p.parseVectorRule()
	default:
		return nil, &MalformedRuleError{}
	}
}
