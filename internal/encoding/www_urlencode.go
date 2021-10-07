package encoding

import (
	"fmt"
	"net/url"
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/helper"
	"openapi-code-sample-generator/internal/log"
	"strings"
)

// UrlencodeParameter encodes the given parameter and its value to application/x-www-form-urlencoded
func UrlencodeParameter(name string, value interface{}) (string, error) {
	encoded := strings.Builder{}

	if helper.IsSlice(value) {
		newVal, ok := value.([]interface{})
		if !ok {
			log.Info(fmt.Sprintf("Type assertion as slice failed for parameter %s with value: %s", name, value))
			return "", errors.TypeAssertionFailed
		}

		encoded.WriteString(url.QueryEscape(name))
		deeperLevels, err := urlencodeSecondLevelObject(newVal)
		if err != nil {
			log.Info(fmt.Sprintf("Can not generate object %s: %s", name, err.Error()))
			return "", err
		}
		encoded.WriteString(deeperLevels)
	} else if helper.IsMap(value) {
		newVal, ok := value.(map[interface{}]interface{})
		if !ok {
			log.Info(fmt.Sprintf("Type assertion as map failed for parameter %s with value: %s", name, value))
			return "", errors.TypeAssertionFailed
		}

		encoded.WriteString(url.QueryEscape(name))
		deeperLevels, err := urlencodeSecondLevelObject(newVal)
		if err != nil {
			log.Info(fmt.Sprintf("Can not generate object %s: %s", name, err.Error()))
			return "", err
		}
		encoded.WriteString(deeperLevels)
	} else {
		newValue := url.QueryEscape(fmt.Sprint(value))
		if stringValue, ok := value.(string); ok {
			if skipParse(stringValue) {
				newValue = stringValue
			}
		}

		encoded.WriteString(fmt.Sprintf("%s=%s", url.QueryEscape(name), newValue))
	}

	return encoded.String(), nil
}

// UrlencodeValue encodes a single value to application/x-www-form-urlencoded
func UrlencodeValue(value interface{}) (string, error) {
	if newValue, ok := value.(string); ok {
		if skipParse(newValue) {
			return newValue, nil
		}
	}
	return url.QueryEscape(fmt.Sprint(value)), nil
}

func urlencodeSecondLevelObject(value interface{}) (string, error) {
	encoded := strings.Builder{}

	if helper.IsMap(value) {
		newVal, ok := value.(map[interface{}]interface{})
		if !ok {
			return "", errors.TypeAssertionFailed
		}
		for key, val := range newVal {
			encoded.WriteString("%5B")
			encoded.WriteString(url.QueryEscape(fmt.Sprint(key)))
			encoded.WriteString("%5D")
			if !helper.IsMap(val) && !helper.IsSlice(newVal) {
				encoded.WriteString("=")
				encoded.WriteString(url.QueryEscape(fmt.Sprint(val)))
			} else {
				va, err := urlencodeSecondLevelObject(val)
				if err != nil {
					return "", err
				}
				encoded.WriteString(fmt.Sprintf(va))
			}
		}
	} else if helper.IsSlice(value) {
		newVal, ok := value.([]interface{})
		if !ok {
			return "", errors.TypeAssertionFailed
		}
		for _, val := range newVal {
			encoded.WriteString("%5B%5D")
			if !helper.IsMap(val) && !helper.IsSlice(newVal) {
				encoded.WriteString("=")
				encoded.WriteString(url.QueryEscape(fmt.Sprint(val)))
			} else {
				va, err := urlencodeSecondLevelObject(val)
				if err != nil {
					return "", err
				}
				encoded.WriteString(fmt.Sprintf(va))
			}
		}
	} else {
		encoded.WriteString("=")
		encoded.WriteString(url.QueryEscape(fmt.Sprint(value)))
	}

	return encoded.String(), nil
}
