package encoding

import (
	"encoding/json"
	"fmt"
)

// JSONEncode groups urlencoding features
type JSONEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to application/x-www-form-urlencoded
func (u *JSONEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	if newValue, ok := value.(map[interface{}]interface{}); ok {
		newMap := make(map[string]interface{})
		for key, val := range newValue {
			newMap[fmt.Sprint(key)] = val
		}
		val, err := json.Marshal(map[string]interface{}{name: newMap})
		return string(val), err
	}

	val, err := json.Marshal(map[string]interface{}{name: value})
	return string(val), err
}

// EnocdeValue encodes a single value to application/x-www-form-urlencoded
func (u *JSONEncode) EnocdeValue(value interface{}) (string, error) {
	if newValue, ok := value.(map[interface{}]interface{}); ok {
		newMap := make(map[string]interface{})
		for key, val := range newValue {
			newMap[fmt.Sprint(key)] = val
		}
		val, err := json.Marshal(newMap)
		return string(val), err
	}

	val, err := json.Marshal(value)
	return string(val), err
}
