package codesample

import (
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/log"

	"github.com/getkin/kin-openapi/openapi3"
)

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

// Constructor holds all information for making samples
type Constructor struct {
	document *openapi3.T
	debug    bool
}

// NewConstructor returns a new constructor
func NewConstructor(document *openapi3.T, debug bool) *Constructor {
	return &Constructor{
		document: document,
		debug:    debug,
	}
}

// AddSamples adds all samples to the document
func (c *Constructor) AddSamples(languages []Language) error {

	for path, pathItem := range c.document.Paths {
		for method, operation := range pathItem.Operations() {
			c.logDebug("# PATH: " + path + " with method: " + method)
			if operation != nil {
				operation.ExtensionProps.Extensions = make(map[string]interface{})
				samples := make([]*CodeSample, 0)

				for _, lang := range languages {
					c.logDebug("## LANGUAGE: " + string(lang))
					sample, err := c.getSample(lang, path, pathItem, operation)
					if err != nil {
						log.Warn(err.Error())
						continue
					}
					samples = append(samples, sample)
				}

				operation.ExtensionProps.Extensions["x-codeSamples"] = samples
			}
		}
	}

	return nil
}

func (c *Constructor) getSample(lang Language, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) (*CodeSample, error) {
	switch lang {
	case LanguageCurl:
		return c.getCurlSample(path, operation, pathItem)
	}

	return nil, errors.UnknownLanguage
}

func (c *Constructor) logDebug(text string) {
	if c.debug {
		log.Debug(text)
	}
}
