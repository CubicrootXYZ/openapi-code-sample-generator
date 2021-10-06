package languages

import (
	"openapi-code-sample-generator/internal/encoding"
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
	cmd.WriteString("?")
	cmd.WriteString(c.getQueryParams(queryParams))
	cmd.WriteString("\" ")

	return &types.CodeSample{
		Lang:   types.LanguageCurl,
		Label:  "curl",
		Source: cmd.String(),
	}, nil
}

func (c *Curl) getQueryParams(params []*types.Parameter) string {
	query := strings.Builder{}
	for i, param := range params {
		encoded, err := encoding.UrlencodeInterface(param.Name, param.Value)
		if err != nil {
			continue
		}

		if i != 0 {
			query.WriteString("&")
		}

		query.WriteString(encoded)
	}

	return query.String()
}

func (c *Curl) getRequestBody() string {
	return ""
}
