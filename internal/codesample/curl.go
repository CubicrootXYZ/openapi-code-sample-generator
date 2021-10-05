package codesample

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func (c *Constructor) getCurlSample(path string, operation *openapi3.Operation, pathItem *openapi3.PathItem) (*CodeSample, error) {
	cmd := strings.Builder{}
	pathParams, queryParams, _, _, err := c.getParameters(operation.Parameters)
	if err != nil {
		return nil, err
	}

	cmd.WriteString("curl \"")
	cmd.WriteString(c.getURL(operation, pathItem))
	cmd.WriteString(getPath(path, pathParams))
	cmd.WriteString("\" ")
	cmd.WriteString(c.getCurlParams(queryParams))

	return &CodeSample{
		Lang:   LanguageCurl,
		Label:  "curl",
		Source: cmd.String(),
	}, nil
}

func (c *Constructor) getCurlParams(params []*parameter) string {
	cmd := strings.Builder{}
	cmd.WriteString("-d \"")

	for i, param := range params {
		if i != 0 {
			cmd.WriteString("&")
		}
		// TODO handle arrays and nested params - move all to "get curl representation or so"
		cmd.WriteString(fmt.Sprintf("%s=%s", url.QueryEscape(param.Name), url.QueryEscape(fmt.Sprint(param.Value))))
	}

	cmd.WriteString("\"")

	return cmd.String()
}
