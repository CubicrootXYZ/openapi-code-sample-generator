package codesample

import (
	"fmt"
	"strings"
)

func getPath(path string, params []*parameter) string {
	for _, param := range params {
		if param == nil {
			continue
		}

		path = strings.Replace(path, fmt.Sprintf("{%s}", param.Name), fmt.Sprint(param.Value), -1)
	}

	return path
}
