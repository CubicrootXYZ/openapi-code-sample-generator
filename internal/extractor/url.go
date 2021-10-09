package extractor

import (
	"openapi-code-sample-generator/internal/errors"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetURL returns the url with multiple fallback strategies
func (o *openAPIExtractor) GetURL(operation *openapi3.Operation, pathItem *openapi3.PathItem, document *openapi3.T) string {
	url, err := o.getURLfromOperation(operation)
	if err != nil {
		url, err := o.getURLfromPath(pathItem)
		if err != nil {
			url, err := o.getURLfromDocument(document)
			if err != nil {
				return "domain.tld"
			}
			return url
		}
		return url
	}

	return url
}

func (o *openAPIExtractor) getURLfromOperation(operation *openapi3.Operation) (string, error) {
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

func (o *openAPIExtractor) getURLfromPath(pathItem *openapi3.PathItem) (string, error) {
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

func (o *openAPIExtractor) getURLfromDocument(document *openapi3.T) (string, error) {
	if document.Servers == nil {
		return "", errors.NoServer
	}

	for _, server := range document.Servers {
		if len(server.URL) > 0 {
			return strings.TrimSuffix(server.URL, "/"), nil
		}
	}

	return "", errors.NoServer
}
