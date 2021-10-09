package extractor

import (
	"fmt"
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/helper"
	"openapi-code-sample-generator/internal/log"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetExampleValueForSchema returns an example value for the given schema
func (o *openAPIExtractor) GetExampleValueForSchema(schema *openapi3.Schema) (interface{}, error) {
	if !helper.IsNil(schema.Example) {
		log.Debug("example")
		return schema.Example, nil
	}

	if !helper.IsNil(schema.Default) {
		log.Debug("default")
		return schema.Default, nil
	}

	for _, value := range schema.Enum {
		if !helper.IsNil(value) {
			log.Debug("enum")
			return value, nil
		}
	}

	val, err := o.getExampleValueByType(schema)
	if err == nil {
		return val, nil
	}

	return nil, errors.UnknownSchema
}

func (o *openAPIExtractor) getExampleValueByType(schema *openapi3.Schema) (interface{}, error) {
	log.Debug(fmt.Sprintf("Schema of type '%s' and format '%s' received", schema.Type, schema.Format))
	switch schema.Type {
	case "integer":
		return 1234, nil
	case "number":
		return 1.234, nil
	case "string":
		switch schema.Format {
		case "byte":
			return []byte("example string"), nil
		case "binary":
			return "01000101 01111000 01100001 01101101 01110000 01101100 01100101", nil
		case "date":
			return time.Now().Format("2006-01-02"), nil
		case "date-time":
			return time.Now().Format(time.RFC3339), nil
		default:
			return "example-string", nil
		}
	case "boolean":
		return false, nil
	case "array":
		if schema.Items != nil && schema.Items.Value != nil {
			val, err := o.GetExampleValueForSchema(schema.Items.Value)
			if err == nil {
				return []interface{}{val}, nil
			}
		}

		return []interface{}{}, nil
	case "object":
		if schema.Properties != nil {
			values := make(map[interface{}]interface{})
			for name, val := range schema.Properties {
				if val == nil || val.Value == nil {
					continue
				}
				newVal, err := o.GetExampleValueForSchema(val.Value)
				if err != nil {
					continue
				}

				values[name] = newVal
			}

			return values, nil
		}

		return []interface{}{}, nil
	}

	log.Warn(fmt.Sprintf("Schema of type '%s' and format '%s' unknown", schema.Type, schema.Format))
	return nil, errors.UnknownSchema
}
