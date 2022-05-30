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

var ReferenceTime = time.Date(2022, 01, 01, 15, 00, 14, 100, time.UTC)

// GetExampleValueForSchema returns an example value for the given schema
func (o *openAPIExtractor) GetExampleValueForSchema(schema *openapi3.Schema, format string) (interface{}, error) {
	log.Debug(fmt.Sprintf("New schema (title: %s; format: %s) received to generate examples for", schema.Title, format))

	if !helper.IsNil(schema.Example) {
		log.Debug("Using example value")
		return schema.Example, nil
	}

	if !helper.IsNil(schema.Default) {
		log.Debug("Using default value")
		return schema.Default, nil
	}

	for _, value := range schema.Enum {
		if !helper.IsNil(value) {
			log.Debug("Using enum value")
			return value, nil
		}
	}

	val, err := o.getExampleValueByType(schema, format)
	if err == nil {
		return val, nil
	}

	log.Debug("Was not able to generate example for schema")
	return nil, errors.ErrUnknownSchema
}

func (o *openAPIExtractor) getExampleValueByType(schema *openapi3.Schema, format string) (interface{}, error) {
	log.Debug("Trying to calculate example for schema")
	switch schema.Type {
	case "integer":
		log.Debug("Using integer default")
		return 1234, nil
	case "number":
		log.Debug("Using number default")
		return 1.234, nil
	case "string":
		switch schema.Format {
		case "byte":
			log.Debug("Using byte default")
			return []byte("example string"), nil
		case "binary":
			log.Debug("Using binary default")
			return "01000101 01111000 01100001 01101101 01110000 01101100 01100101", nil
		case "date":
			log.Debug("Using date default")
			return ReferenceTime.Format("2006-01-02"), nil
		case "date-time":
			log.Debug("Using date-time default")
			return ReferenceTime.Format(time.RFC3339), nil
		default:
			log.Debug("Using string default")
			return "example-string", nil
		}
	case "boolean":
		log.Debug("Using boolean default")
		return false, nil
	case "array":
		log.Debug("Generating array example")
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
		log.Debug("Generating object example")
		return o.getExampleValueForObject(schema, format)
	}

	log.Warn(fmt.Sprintf("Schema of type '%s' and format '%s' unknown", schema.Type, schema.Format))
	return nil, errors.ErrUnknownSchema
}

// Objects are way more complex then primitive values
// We need to consider embedding here
func (o *openAPIExtractor) getExampleValueForObject(schema *openapi3.Schema, format string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	if schema.Properties != nil {
		log.Debug("Adding properties to object")
		for name, val := range schema.Properties {
			if val == nil || val.Value == nil {
				log.Debug(fmt.Sprintf("Skipping property %s with empty value", name))
				continue
			}
			newVal, err := o.GetExampleValueForSchema(val.Value, format)
			if err != nil {
				log.Debug(fmt.Sprintf("Skipping property %s with error: %s", name, err.Error()))
				continue
			}

			if xmlInfo, ok := val.Value.XML.(map[string]interface{}); ok && format == types.EncodingXML {
				if value, ok := xmlInfo["name"]; ok {
					log.Debug(fmt.Sprintf("Overriding property %s to: %s due to XML tag given", name, value))
					name = fmt.Sprint(value)
				}
			}

			values[name] = newVal
		}

	}

	// oneOf: pick the first one, 0 is a pointer to the main schema itself
	if schema.OneOf != nil && len(schema.OneOf) > 1 && schema.OneOf[1].Value != nil {
		log.Debug(fmt.Sprintf("Adding first of 'oneOf' embedding (ref: %s)", schema.OneOf[1].Ref))
		additionalValues, err := o.getExampleValueForObject(schema.OneOf[1].Value, format)
		if err == nil {
			for name, value := range additionalValues {
				values[name] = value
			}
		}
	}

	// anyOf: pick the first one, 0 is a pointer to the main schema itself
	if schema.AnyOf != nil && len(schema.AnyOf) > 1 && schema.AnyOf[1].Value != nil {
		log.Debug(fmt.Sprintf("Adding first of 'anyOf' embedding (ref: %s)", schema.AnyOf[1].Ref))
		additionalValues, err := o.getExampleValueForObject(schema.AnyOf[1].Value, format)
		if err == nil {
			for name, value := range additionalValues {
				values[name] = value
			}
		}
	}

	// allOf: add all schemas, 0 is a pointer to the main schema itself
	if schema.AllOf != nil {
		log.Debug("Adding all of 'allOf' embedding")
		for i, additionalSchema := range schema.AllOf {
			if additionalSchema.Value == nil || i == 0 {
				continue
			}

			log.Debug(fmt.Sprintf("Adding %d of 'allOf' embedding (ref: %s)", i, additionalSchema.Ref))

			additionalValues, err := o.getExampleValueForObject(additionalSchema.Value, format)
			if err == nil {
				for name, value := range additionalValues {
					values[name] = value
				}
			}
		}
	}

	return values, nil
}
