package types

import "github.com/getkin/kin-openapi/openapi3"

// Generator is an interface for language specific generators to adopt
type Generator interface {
	GetSample(path string, operation *openapi3.Operation, pathItem *openapi3.PathItem) (*CodeSample, error)
}
