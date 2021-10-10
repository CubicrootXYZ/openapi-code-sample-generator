package encoding

import (
	"fmt"
	"strings"

	"github.com/clbanning/anyxml"
)

// XMLEncode groups xml encoding features
type XMLEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to application/x-www-form-urlencoded
func (x *XMLEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	if newValue, ok := value.(map[interface{}]interface{}); ok {
		newMap := make(map[string]interface{})
		for key, val := range newValue {
			newMap[fmt.Sprint(key)] = val
		}
		val, err := anyxml.Xml(map[string]interface{}{name: newMap})
		return string(val), err
	}

	val, err := anyxml.Xml(map[string]interface{}{name: value})
	return string(val), err
}

// EnocdeValue encodes a single value to application/x-www-form-urlencoded
func (x *XMLEncode) EnocdeValue(ref string, value interface{}) (string, error) {
	root := x.rootTag(ref)

	if newValue, ok := value.(map[interface{}]interface{}); ok {
		newMap := make(map[string]interface{})
		for key, val := range newValue {
			newMap[fmt.Sprint(key)] = val
		}
		val, err := anyxml.Xml(newMap, root)
		return XmlProlog() + string(val), err
	}

	val, err := anyxml.Xml(value, root)
	return XmlProlog() + string(val), err
}

func XmlProlog() string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
}

func (x *XMLEncode) rootTag(ref string) string {
	if ref == "" {
		return "doc"
	}

	parts := strings.Split(ref, "/")

	if len(parts) == 0 {
		return "doc"
	}

	return parts[len(parts)-1]
}
