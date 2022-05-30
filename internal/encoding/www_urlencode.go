package encoding

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/helper"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

// URLEncode groups urlencoding features
type URLEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to application/x-www-form-urlencoded
func (u *URLEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	encoded := strings.Builder{}

	if helper.IsSlice(value) {
		newVal, ok := value.([]interface{})
		if !ok {
			log.Info(fmt.Sprintf("Type assertion as slice failed for parameter %s with value: %s", name, value))
			return "", errors.ErrTypeAssertionFailed
		}

		encoded.WriteString(url.QueryEscape(name))
		deeperLevels, err := u.urlencodeSecondLevelObject(newVal)
		if err != nil {
			log.Info(fmt.Sprintf("Can not generate object %s: %s", name, err.Error()))
			return "", err
		}
		encoded.WriteString(deeperLevels)
	} else if helper.IsMap(value) {
		newVal, ok := value.(map[string]interface{})
		if !ok {
			log.Info(fmt.Sprintf("Type assertion as map failed for parameter %s with value: %s", name, value))
			return "", errors.ErrTypeAssertionFailed
		}

		encoded.WriteString(url.QueryEscape(name))
		deeperLevels, err := u.urlencodeSecondLevelObject(newVal)
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

// EnocdeValue encodes a single value to application/x-www-form-urlencoded
func (u *URLEncode) EnocdeValue(ref string, value interface{}, meta *types.FormattingMeta) (string, error) {
	encoded := strings.Builder{}

	if helper.IsSlice(value) {
		newVal, ok := value.([]interface{})
		if !ok {
			log.Info(fmt.Sprintf("Type assertion as slice failed for value: %s", value))
			return "", errors.ErrTypeAssertionFailed
		}

		for index, newVa := range newVal {
			encoded.WriteString(fmt.Sprint(index))

			deeperLevels, err := u.urlencodeSecondLevelObject(newVa)
			if err != nil {
				log.Info(fmt.Sprintf("Can not generate object %s", err.Error()))
				return "", err
			}
			encoded.WriteString(deeperLevels)
		}

	} else if helper.IsMap(value) {
		newVal, ok := value.(map[string]interface{})
		if !ok {
			log.Info(fmt.Sprintf("Type assertion as map failed for value: %s", value))
			return "", errors.ErrTypeAssertionFailed
		}

		i := 0
		for key, newVa := range newVal {
			if i != 0 {
				encoded.WriteString("&")
			}

			encoded.WriteString(fmt.Sprint(key))
			deeperLevels, err := u.urlencodeSecondLevelObject(newVa)
			if err != nil {
				log.Info(fmt.Sprintf("Can not generate object: %s", err.Error()))
				return "", err
			}

			encoded.WriteString(deeperLevels)

			i++
		}

	} else {
		newValue := url.QueryEscape(fmt.Sprint(value))
		if stringValue, ok := value.(string); ok {
			if skipParse(stringValue) {
				newValue = stringValue
			}
		}

		encoded.WriteString(newValue)
	}

	return encoded.String(), nil
}

func (u *URLEncode) urlencodeSecondLevelObject(value interface{}) (string, error) {
	encoded := strings.Builder{}

	if helper.IsMap(value) {
		newVal, ok := value.(map[string]interface{})
		if !ok {
			log.Debug("IsMap but is not interface map")
			return "", errors.ErrTypeAssertionFailed
		}
		for key, val := range newVal {
			encoded.WriteString("%5B")
			encoded.WriteString(url.QueryEscape(fmt.Sprint(key)))
			encoded.WriteString("%5D")
			if !helper.IsMap(val) && !helper.IsSlice(val) {
				encoded.WriteString("=")
				encoded.WriteString(url.QueryEscape(fmt.Sprint(val)))
			} else {
				va, err := u.urlencodeSecondLevelObject(val)
				if err != nil {
					return "", err
				}
				encoded.WriteString(va)
			}
		}
	} else if helper.IsSlice(value) {
		newVal, ok := value.([]interface{})
		if !ok {
			log.Debug("IsSlice but is not interface slice")
			return "", errors.ErrTypeAssertionFailed
		}
		for _, val := range newVal {
			encoded.WriteString("%5B%5D")
			if !helper.IsMap(val) && !helper.IsSlice(val) {
				encoded.WriteString("=")
				encoded.WriteString(url.QueryEscape(fmt.Sprint(val)))
			} else {
				va, err := u.urlencodeSecondLevelObject(val)
				if err != nil {
					return "", err
				}
				encoded.WriteString(va)
			}
		}
	} else {
		encoded.WriteString("=")
		encoded.WriteString(url.QueryEscape(fmt.Sprint(value)))
	}

	return encoded.String(), nil
}
