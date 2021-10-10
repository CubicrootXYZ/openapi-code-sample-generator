package extractor

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/helper"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetExampleValueForSchema returns an example value for the given schema
func (o *openAPIExtractor) GetExampleValueForSchema(schema *openapi3.Schema, format string) (interface{}, error) {
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

	val, err := o.getExampleValueByType(schema, format)
	if err == nil {
		return val, nil
	}

	return nil, errors.UnknownSchema
}

func (o *openAPIExtractor) getExampleValueByType(schema *openapi3.Schema, format string) (interface{}, error) {
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
			val, err := o.GetExampleValueForSchema(schema.Items.Value, format)

			// Handle wrapping
			if xmlInfo, ok := schema.XML.(map[string]interface{}); ok && format == types.EncodingXML {
				if wrapped, ok := xmlInfo["wrapped"]; ok {
					if fmt.Sprint(wrapped) == "true" || fmt.Sprint(wrapped) == "True" {
						// Check child for tag names
						if xmlChildInfo, ok := schema.Items.Value.XML.(map[string]interface{}); ok && format == types.EncodingXML {
							if wrappedName, ok := xmlChildInfo["name"]; ok {
								log.Debug(fmt.Sprintf("Wrapping with tag names: %s", fmt.Sprint(wrappedName)))
								return map[string]interface{}{fmt.Sprint(wrappedName): val}, nil
							}
						}
					}
				}
			}

			if err == nil {
				return []interface{}{val}, nil
			}
		}

		return []interface{}{}, nil
	case "object":
		if schema.Properties != nil {
			values := make(map[string]interface{})
			for name, val := range schema.Properties {
				if val == nil || val.Value == nil {
					log.Debug(fmt.Sprintf("Skipping %s with empty value", name))
					continue
				}
				newVal, err := o.GetExampleValueForSchema(val.Value, format)
				if err != nil {
					log.Debug(fmt.Sprintf("Skipping %s with: %s", name, err.Error()))
					continue
				}

				if xmlInfo, ok := val.Value.XML.(map[string]interface{}); ok && format == types.EncodingXML {
					if value, ok := xmlInfo["name"]; ok {
						log.Debug(fmt.Sprintf("Overriding xml tag %s to: %s", name, value))
						name = fmt.Sprint(value)
					}
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
