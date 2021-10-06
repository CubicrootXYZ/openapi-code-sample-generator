package helper

import (
	"fmt"
	"openapi-code-sample-generator/internal/types"
	"strings"
)

// GetPath returns the path with sample params set
func GetPath(path string, params []*types.Parameter) string {
	for _, param := range params {
		if param == nil {
			continue
		}

		path = strings.Replace(path, fmt.Sprintf("{%s}", param.Name), fmt.Sprint(param.Value), -1)
	}

	return path
}
