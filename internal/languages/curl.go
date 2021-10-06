package languages

import (
	"fmt"
	"net/url"
	"openapi-code-sample-generator/internal/helper"
	"openapi-code-sample-generator/internal/types"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Curl holds information for curl samples
type Curl struct {
	document *openapi3.T
}

// NewCurl returns a new curl object
func NewCurl(document *openapi3.T) types.Generator {
	return &Curl{
		document: document,
	}
}

// GetSample returns a curl sample for the given operation
func (c *Curl) GetSample(path string, operation *openapi3.Operation, pathItem *openapi3.PathItem) (*types.CodeSample, error) {
	cmd := strings.Builder{}
	pathParams, queryParams, _, _, err := helper.GetParameters(operation.Parameters)
	if err != nil {
		return nil, err
	}

	cmd.WriteString("curl \"")
	cmd.WriteString(helper.GetURL(operation, pathItem, c.document))
	cmd.WriteString(helper.GetPath(path, pathParams))
	cmd.WriteString("\" ")
	cmd.WriteString(c.getCurlParams(queryParams))

	return &types.CodeSample{
		Lang:   types.LanguageCurl,
		Label:  "curl",
		Source: cmd.String(),
	}, nil
}

func (c *Curl) getCurlParams(params []*types.Parameter) string {
	cmd := strings.Builder{}
	cmd.WriteString("-d \"")

	for i, param := range params {
		if i != 0 {
			cmd.WriteString("&")
		}

		cmd.WriteString(fmt.Sprintf("%s=%s", url.QueryEscape(param.Name), url.QueryEscape(fmt.Sprint(param.Value))))
	}

	cmd.WriteString("\"")

	return cmd.String()
}
