package types

import "github.com/getkin/kin-openapi/openapi3"

// Extractor defines an interface for feature extractors
type Extractor interface {
	GetParameters(params openapi3.Parameters) (Parameters, error)
	GetPathExample(path string, params []*Parameter) string
	GetURL(operation *openapi3.Operation, pathItem *openapi3.PathItem, document *openapi3.T) string
	GetRequestBody(body *openapi3.RequestBody) (value interface{}, format string, err error)
	GetExampleValueForSchema(schema *openapi3.Schema, format string) (interface{}, error)
	GetSecurity(operation *openapi3.Operation, document *openapi3.T) (params Parameters, basicAuth bool, err error)
}
