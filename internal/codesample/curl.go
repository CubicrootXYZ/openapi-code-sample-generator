package codesample

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func (c *Constructor) getCurlSample(path string, operation *openapi3.Operation, pathItem *openapi3.PathItem) (*CodeSample, error) {
	cmd := strings.Builder{}
	pathParams, _, _, _, err := c.getParameters(operation.Parameters)
	if err != nil {
		return nil, err
	}

	cmd.WriteString("curl \"")
	cmd.WriteString(c.getURL(operation, pathItem))
	cmd.WriteString(getPath(path, pathParams))
	cmd.WriteString("\"")

	return &CodeSample{
		Lang:   LanguageCurl,
		Label:  "curl",
		Source: cmd.String(),
	}, nil
}
