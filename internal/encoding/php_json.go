package encoding

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

// PhpJsonEncode groups json features for php
type PhpJsonEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to json
func (j *PhpJsonEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	return j.EnocdeValue("", map[string]interface{}{name: value}, nil)
}

// EnocdeValue encodes a single value to a php json object
func (j *PhpJsonEncode) EnocdeValue(ref string, value interface{}, meta *types.FormattingMeta) (string, error) {
	switch value.(type) {
	case string:
		return `"` + value.(string) + `"`, nil
	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64, float64, float32:
		return fmt.Sprintf(`"%v"`, value), nil
	case bool:
		return fmt.Sprint(value), nil
	case byte:
		return `"` + string(value.(byte)) + `"`, nil
	case []interface{}:
		out := strings.Builder{}
		out.WriteString("json_encode(array(")

		for _, item := range value.([]interface{}) {
			itemStringified, err := j.EnocdeValue(ref, item, meta)
			if err != nil {
				log.Warn(err.Error())
				continue
			}

			out.WriteString(itemStringified)
			out.WriteString(",")
		}
		out.WriteString("))")

		return out.String(), nil
	case map[string]interface{}:
		out := strings.Builder{}
		out.WriteString("json_encode(array(")

		for key, item := range value.(map[string]interface{}) {
			itemStringified, err := j.EnocdeValue(ref, item, meta)
			if err != nil {
				log.Warn(err.Error())
				continue
			}

			out.WriteString("\"")
			out.WriteString(key)
			out.WriteString("\" => ")
			out.WriteString(itemStringified)
			out.WriteString(",")
		}
		out.WriteString("))")

		return out.String(), nil
	default:
		log.Info(fmt.Sprintf("Type %T is not known for PHP-JSON encoder", value))
		return "", errors.ErrUnknownType
	}
}
