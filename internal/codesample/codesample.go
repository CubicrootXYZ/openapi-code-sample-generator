package codesample

import (
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/log"
	"openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// Orchestrator orchestrates the sample generation
type Orchestrator struct {
	document   *openapi3.T
	generators map[types.Language]types.Generator
}

// NewConstructor returns a new constructor
func NewConstructor(document *openapi3.T, generators map[types.Language]types.Generator) *Orchestrator {
	return &Orchestrator{
		document:   document,
		generators: generators,
	}
}

// AddSamples adds all samples to the document
func (o *Orchestrator) AddSamples(languages []types.Language) error {
	// Iterate over paths
	for path, pathItem := range o.document.Paths {
		// Iterate over operations
		for method, operation := range pathItem.Operations() {
			log.Debug("# PATH: " + path + " with method: " + method)

			if operation == nil {
				continue
			}

			operation.ExtensionProps.Extensions = make(map[string]interface{})
			samples, err := o.getSamples(languages, path, pathItem, operation)
			if err != nil {
				log.Warn("Can not generate samples: " + err.Error())
				continue
			}
			operation.ExtensionProps.Extensions["x-codeSamples"] = samples

		}
	}

	return nil
}

func (o *Orchestrator) getSamples(languages []types.Language, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) ([]*types.CodeSample, error) {
	samples := make([]*types.CodeSample, 0)

	for _, lang := range languages {
		log.Debug("## LANGUAGE: " + string(lang))

		sample, err := o.getSample(lang, path, pathItem, operation)
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		samples = append(samples, sample)
	}

	return samples, nil
}

func (o *Orchestrator) getSample(lang types.Language, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) (*types.CodeSample, error) {
	generator, ok := o.generators[lang]
	if !ok {
		return nil, errors.UnknownLanguage
	}

	return generator.GetSample(path, operation, pathItem, o.document)
}
