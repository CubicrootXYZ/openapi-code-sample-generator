package languages

import (
	"openapi-code-sample-generator/internal/types"
)

// Generators returns all available generators
func Generators(encoders map[string]types.Encoder, extractor types.Extractor) map[types.Language]types.Generator {
	generators := make(map[types.Language]types.Generator)
	generators[types.LanguageCurl] = NewCurl(encoders, extractor)
	return generators
}