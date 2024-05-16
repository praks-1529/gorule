// File: parser_utils.go
// Common utility functions
package gorule

import (
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// StringToInterface takes a string and returns the actual type wrapped in interface
func StringToInterface(value string) interface{} {
	switch value {
	case "true":
		return true
	case "false":
		return false
	default:
		if intValue, err := strconv.Atoi(string(value)); err == nil {
			return intValue
		} else if floatValue, err := strconv.ParseFloat(string(value), 64); err == nil {
			return floatValue
		} else {
			return value
		}
	}
}

// resolveValue takes a key and then gets the coresponding value from JSON
// Example if key = a.b and ipData =  { a : { b: 10 }}.
// Then return value is int(10)
func resolveValue(key string, ipData []byte) interface{} {
	value := gjson.Get(string(ipData), key)
	if value.IsObject() {
		// Must be only called for literals
		panic("Expecting object but did not find")
	}
	switch value.Type {
	case gjson.String:
		return value.Raw
	case gjson.Number:
		return StringToInterface(value.Raw)
	case gjson.False:
		return value.Bool()
	case gjson.True:
		return value.Bool()
	case gjson.Null:
		return key
	default:
		return key
	}
}

// resolveLength takes a key in the form of a.size() and then gets the length
// Example if key = a.size() and ipData is a : [{}, {}, {}]
// Then return value is int(3)
func resolveLength(key string, ipData []byte) interface{} {
	value := gjson.Get(string(ipData), key)
	if !(value.IsArray()) {
		panic("Expecting array but not found")
	}
	return len(value.Array())
}

func hasArrayIndex(key string) bool {
	if strings.Contains(key, "[") && strings.Contains(key, "]") {
		return true
	}
	return false
}

func getStringFromToken(token Token) string {
	return string(token)
}
