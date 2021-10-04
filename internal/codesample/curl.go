package codesample

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func getCurlSample(operation *openapi3.Operation, pathItem *openapi3.PathItem) *CodeSample {
	cmd := strings.Builder{}

	cmd.WriteString("curl \"")
	cmd.WriteString(getUURLfromOperation(operation))
	cmd.WriteString(pathItem.Ref)
	cmd.WriteString("\"")

	return &CodeSample{
		Lang:   LanguageCurl,
		Label:  "curl",
		Source: cmd.String(),
	}
}

func getUURLfromOperation(operation *openapi3.Operation) string {
	if operation.Servers == nil {
		return "domain.tld"
	}

	for _, server := range *operation.Servers {
		if len(server.URL) > 0 {
			return strings.TrimSuffix(server.URL, "/")
		}
	}

	return "domain.tld"
}
