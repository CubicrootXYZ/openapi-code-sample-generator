package codesample

import (
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/languages"
	"openapi-code-sample-generator/internal/log"
	"openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// Constructor holds all information for making samples
type Constructor struct {
	document   *openapi3.T
	debug      bool
	generators map[types.Language]types.Generator
}

// NewConstructor returns a new constructor
func NewConstructor(document *openapi3.T, debug bool) *Constructor {
	generators := make(map[types.Language]types.Generator)
	generators[types.LanguageCurl] = languages.NewCurl(document)

	return &Constructor{
		document:   document,
		debug:      debug,
		generators: generators,
	}
}

// AddSamples adds all samples to the document
func (c *Constructor) AddSamples(languages []types.Language) error {
	for path, pathItem := range c.document.Paths {
		for method, operation := range pathItem.Operations() {
			log.Debug("# PATH: " + path + " with method: " + method)
			if operation != nil {
				operation.ExtensionProps.Extensions = make(map[string]interface{})
				samples := make([]*types.CodeSample, 0)

				for _, lang := range languages {
					log.Debug("## LANGUAGE: " + string(lang))
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

func (c *Constructor) getSample(lang types.Language, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) (*types.CodeSample, error) {
	generator, ok := c.generators[lang]
	if !ok {
		return nil, errors.UnknownLanguage
	}

	return generator.GetSample(path, operation, pathItem)
}
