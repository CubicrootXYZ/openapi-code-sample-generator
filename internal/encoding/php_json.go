package encoding

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

// PhpEncode groups json features for php
type PhpEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to json
func (j *PhpEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	return j.EnocdeValue("", map[string]interface{}{name: value}, nil)
}

// EnocdeValue encodes a single value to a php json object
func (j *PhpEncode) EnocdeValue(ref string, value interface{}, meta *types.FormattingMeta) (string, error) {
	switch t := value.(type) {
	case string:
		return `"` + t + `"`, nil
	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64, float64, float32:
		return fmt.Sprintf(`"%v"`, t), nil
	case bool:
		return fmt.Sprint(value), nil
	case byte:
		return `"` + string(t) + `"`, nil
	case []interface{}:
		out := strings.Builder{}
		out.WriteString("array(\n")

		for _, item := range t {
			itemStringified, err := j.EnocdeValue(ref, item, meta)
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			out.WriteString("\t")
			out.WriteString(strings.Replace(itemStringified, "\n", "\n\t", -1))
			out.WriteString(",\n")
		}
		out.WriteString(")")

		return out.String(), nil
	case map[string]interface{}:
		out := strings.Builder{}
		out.WriteString("array(\n")

		for key, item := range t {
			itemStringified, err := j.EnocdeValue(ref, item, meta)
			if err != nil {
				log.Warn(err.Error())
				continue
			}

			out.WriteString("\t\"")
			out.WriteString(key)
			out.WriteString("\" => ")
			out.WriteString(strings.Replace(itemStringified, "\n", "\n\t", -1))
			out.WriteString(",\n")
		}
		out.WriteString(")")

		return out.String(), nil
	default:
		log.Info(fmt.Sprintf("Type %T is not known for PHP-JSON encoder", value))
		return "", errors.ErrUnknownType
	}
}
