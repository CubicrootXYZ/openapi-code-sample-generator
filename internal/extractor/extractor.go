package extractor

import "github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

type openAPIExtractor struct {
}

// NewOpenAPIExtractor returns a new extractor for openapi specs
func NewOpenAPIExtractor() types.Extractor {
	return &openAPIExtractor{}
}
