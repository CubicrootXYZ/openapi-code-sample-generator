package codesample

import (
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// Executor orchestrates the sample generation
type Executor struct {
	document  *openapi3.T
	templater templater.Templater
}

// NewExecutor returns a new constructor
func NewExecutor(document *openapi3.T, templater templater.Templater) *Executor {
	return &Executor{
		document:  document,
		templater: templater,
	}
}

// AddSamples adds all samples to the document
func (o *Executor) AddSamples(languages []types.Language) error {
	// Iterate over paths
	for path, pathItem := range o.document.Paths.Map() {
		// Iterate over operations
		for method, operation := range pathItem.Operations() {
			log.Debug("# PATH: " + path + " with method: " + method)

			if operation == nil {
				continue
			}

			operation.Extensions = make(map[string]interface{})
			samples, err := o.getSamples(languages, method, path, pathItem, operation)
			if err != nil {
				log.Warn("Can not generate samples: " + err.Error())
				continue
			}
			operation.Extensions["x-codeSamples"] = samples

		}
	}

	return nil
}

func (o *Executor) getSamples(languages []types.Language, httpVerb, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) ([]*types.CodeSample, error) {
	samples := make([]*types.CodeSample, 0)

	for _, lang := range languages {
		log.Debug("## LANGUAGE: " + string(lang))

		sample, err := o.getTemplatedSample(lang, httpVerb, path, pathItem, operation)
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		samples = append(samples, sample)
	}

	return samples, nil
}

func (o *Executor) getTemplatedSample(lang types.Language, httpVerb, path string, pathItem *openapi3.PathItem, operation *openapi3.Operation) (*types.CodeSample, error) {
	return o.templater.Template(lang, templater.NewEndpoint(httpVerb, path, operation, pathItem, o.document))
}
