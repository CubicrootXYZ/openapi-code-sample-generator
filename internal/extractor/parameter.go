package extractor

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/helper"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetParameters returns the different parameter types
func (o *openAPIExtractor) GetParameters(params openapi3.Parameters) (types.Parameters, error) {
	parameters := types.Parameters{}
	parameters.Path = make([]*types.Parameter, 0)
	parameters.Query = make([]*types.Parameter, 0)
	parameters.Header = make([]*types.Parameter, 0)
	parameters.Cookie = make([]*types.Parameter, 0)

	if params == nil {
		return parameters, nil
	}

	for _, ref := range params {
		if ref == nil || ref.Value == nil {
			continue
		}

		// Only use required parameters
		if ref.Value.Required && !ref.Value.Deprecated {
			log.Debug("### Param " + ref.Value.Name + " in " + ref.Value.In)
			val, err := o.getParamValue(ref.Value)
			if err != nil {
				return parameters, errors.ErrUnknownParameter
			}

			log.Debug(fmt.Sprint("is set to ", val))

			param := &types.Parameter{
				Name:  ref.Value.Name,
				Value: val,
			}

			switch strings.ToLower(ref.Value.In) {
			case "path":
				parameters.Path = append(parameters.Path, param)
			case "query":
				parameters.Query = append(parameters.Query, param)
			case "head", "header":
				parameters.Header = append(parameters.Header, param)
			case "cookie":
				parameters.Cookie = append(parameters.Cookie, param)
			}
		}
	}

	return parameters, nil
}

func (o *openAPIExtractor) getParamValue(param *openapi3.Parameter) (interface{}, error) {
	if !helper.IsNil(param.Example) {
		log.Debug("using param example value")
		return param.Example, nil
	}

	if param.Schema != nil && param.Schema.Value != nil {
		val, err := o.GetExampleValueForSchema(param.Schema.Value, "")
		if err == nil {
			return val, nil
		}
	}

	return nil, errors.ErrUnknownParameter
}
