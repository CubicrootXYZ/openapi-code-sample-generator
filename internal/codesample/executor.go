package codesample

import (
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// Executor orchestrates the sample generation
type Executor struct {
	document   *openapi3.T
	generators map[types.Language]types.Generator
}

// NewExecutor returns a new constructor
func NewExecutor(document *openapi3.T, generators map[types.Language]types.Generator) *Executor {
	return &Executor{
		document:   document,
		generators: generators,
	}
}

// AddSamples adds all samples to the document
func (o *Executor) AddSamples(languages []types.Language) error {
	// Iterate over paths
	for path, pathItem := range o.document.Paths {
		// Iterate over operations
		for method, operation := range pathItem.Operations() {
			log.Debug("# PATH: " + path + " with method: " + method)

			if operation == nil {
				continue
			}

			operation.ExtensionProps.Extensions = make(map[string]interface{})
			samples, err := o.getSamples(languages, method, path, pathItem, operation)
			if err != nil {
				log.Warn("Can not generate samples: " + err.Error())
				continue
			}
			operation.ExtensionProps.Extensions["x-codeSamples"] = samples

		}
	}

	return nil
}

func (o *Executor) getSamples(languages []types.Language, httpVerb, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) ([]*types.CodeSample, error) {
	samples := make([]*types.CodeSample, 0)

	for _, lang := range languages {
		log.Debug("## LANGUAGE: " + string(lang))

		sample, err := o.getSample(lang, httpVerb, path, pathItem, operation)
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		samples = append(samples, sample)
	}

	return samples, nil
}

func (o *Executor) getSample(lang types.Language, httpVerb, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) (*types.CodeSample, error) {
	generator, ok := o.generators[lang]
	if !ok {
		return nil, errors.UnknownLanguage
	}

	return generator.GetSample(httpVerb, path, operation, pathItem, o.document)
}
