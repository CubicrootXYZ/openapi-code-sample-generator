package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

// JSONEncode groups json features
type JSONEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to json
func (j *JSONEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	return j.EnocdeValue("", map[string]interface{}{name: value}, nil)
}

// EnocdeValue encodes a single value to json
func (j *JSONEncode) EnocdeValue(ref string, value interface{}, meta *types.FormattingMeta) (string, error) {
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
