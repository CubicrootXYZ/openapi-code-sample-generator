package helper

import (
	"fmt"
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/log"
	"openapi-code-sample-generator/internal/types"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetParameters returns the different parameter types
func GetParameters(params openapi3.Parameters) (pathParams []*types.Parameter, queryParams []*types.Parameter, headerParams []*types.Parameter, cookieParams []*types.Parameter, err error) {
	err = nil
	pathParams = make([]*types.Parameter, 0)
	queryParams = make([]*types.Parameter, 0)
	headerParams = make([]*types.Parameter, 0)
	cookieParams = make([]*types.Parameter, 0)

	if params == nil {
		return
	}

	for _, ref := range params {
		if ref == nil || ref.Value == nil {
			continue
		}

		// Only use required parameters
		if ref.Value.Required && !ref.Value.Deprecated {
			log.Debug("### Param " + ref.Value.Name + " in " + ref.Value.In)
			val, err := getParamValue(ref.Value)
			if err != nil {
				return nil, nil, nil, nil, errors.UnknownParameter
			}

			log.Debug(fmt.Sprint("is set to ", val))

			param := &types.Parameter{
				Name:  ref.Value.Name,
				Value: val,
			}

			switch strings.ToLower(ref.Value.In) {
			case "path":
				pathParams = append(pathParams, param)
			case "query":
				queryParams = append(queryParams, param)
			case "head", "header":
				headerParams = append(headerParams, param)
			case "cookie":
				cookieParams = append(cookieParams, param)
			}
		}
	}

	return
}

func getParamValue(param *openapi3.Parameter) (interface{}, error) {
	if !IsNil(param.Example) {
		log.Debug("using param example value")
		return param.Example, nil
	}

	if param.Schema != nil && param.Schema.Value != nil {
		val, err := GetExampleValueForSchema(param.Schema.Value)
		if err == nil {
			return val, nil
		}
	}

	return nil, errors.UnknownParameter
}
