package extractor

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

// GetPath returns the path with sample params set
func (o *openAPIExtractor) GetPathExample(path string, params []*types.Parameter) string {
	for _, param := range params {
		if param == nil {
			continue
		}

		path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", param.Name), fmt.Sprint(param.Value))
	}

	return path
}
