package codesample

import (
	"openapi-code-sample-generator/internal/errors"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func (c *Constructor) getURL(operation *openapi3.Operation, pathItem *openapi3.PathItem) string {
	url, err := c.getURLfromOperation(operation)
	if err != nil {
		url, err := c.getURLfromPath(pathItem)
		if err != nil {
			url, err := c.getURLfromDocument()
			if err != nil {
				return "domain.tld"
			}
			return url
		}
		return url
	}

	return url
}

func (c *Constructor) getURLfromOperation(operation *openapi3.Operation) (string, error) {
	if operation.Servers == nil {
		return "", errors.NoServer
	}

	for _, server := range *operation.Servers {
		if len(server.URL) > 0 {
			return strings.TrimSuffix(server.URL, "/"), nil
		}
	}

	return "", errors.NoServer
}

func (c *Constructor) getURLfromPath(pathItem *openapi3.PathItem) (string, error) {
	if pathItem.Servers == nil {
		return "", errors.NoServer
	}

	for _, server := range pathItem.Servers {
		if len(server.URL) > 0 {
			return strings.TrimSuffix(server.URL, "/"), nil
		}
	}

	return "", errors.NoServer
}

func (c *Constructor) getURLfromDocument() (string, error) {
	if c.document.Servers == nil {
		return "", errors.NoServer
	}

	for _, server := range c.document.Servers {
		if len(server.URL) > 0 {
			return strings.TrimSuffix(server.URL, "/"), nil
		}
	}

	return "", errors.NoServer
}
