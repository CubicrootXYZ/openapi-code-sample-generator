package codesample

import "github.com/getkin/kin-openapi/openapi3"

// CodeSample represents a single code example
type CodeSample struct {
	Lang   Language `yaml:"lang" json:"lang"`     // language of the sample
	Source string   `yaml:"source" json:"source"` // the actual source code
	Label  string   `yaml:"label" json:"label"`   // displayed language name
}

// Language enum
type Language string

const (
	// LanguageCurl curl
	LanguageCurl = Language("curl")
)

// GetSamples returns a single sample for the given language
func GetSamples(lang Language, operation *openapi3.Operation, pathItem *openapi3.PathItem) *[]*CodeSample {
	samples := make([]*CodeSample, 0)
	switch lang {
	case LanguageCurl:
		samples = append(samples, getCurlSample(operation, pathItem))
	}

	return &samples
}
