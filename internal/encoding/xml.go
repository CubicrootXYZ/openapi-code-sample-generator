package encoding

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
	"github.com/clbanning/anyxml"
)

// XMLEncode groups xml encoding features
type XMLEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to xml
func (x *XMLEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	return x.EnocdeValue("", map[string]interface{}{name: value}, nil)
}

// EnocdeValue encodes a single value to xml
func (x *XMLEncode) EnocdeValue(ref string, value interface{}, meta *types.FormattingMeta) (string, error) {
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
