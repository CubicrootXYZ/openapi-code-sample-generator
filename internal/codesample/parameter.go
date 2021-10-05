package codesample

import (
	"fmt"
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/helper"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type parameter struct {
	Name  string
	Value interface{}
}

func (c *Constructor) getParameters(params openapi3.Parameters) (pathParams []*parameter, queryParams []*parameter, headerParams []*parameter, cookieParams []*parameter, err error) {
	err = nil
	pathParams = make([]*parameter, 0)
	queryParams = make([]*parameter, 0)
	headerParams = make([]*parameter, 0)
	cookieParams = make([]*parameter, 0)

	if params == nil {
		return
	}

	for _, ref := range params {
		if ref == nil || ref.Value == nil {
			continue
		}

		// Only use required parameters
		if ref.Value.Required && !ref.Value.Deprecated {
			c.logDebug("### Param " + ref.Value.Name + " in " + ref.Value.In)
			val, err := c.getParamValue(ref.Value)
			if err != nil {
				return nil, nil, nil, nil, errors.UnknownParameter
			}

			c.logDebug(fmt.Sprint("is set to ", val))

			param := &parameter{
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

func (c *Constructor) getParamValue(param *openapi3.Parameter) (interface{}, error) {
	if !helper.IsNil(param.Example) {
		c.logDebug("using param example value")
		return param.Example, nil
	}

	if param.Schema != nil && param.Schema.Value != nil {
		val, err := helper.GetExampleValueForSchema(param.Schema.Value)
		if err == nil {
			return val, nil
		}
	}

	return nil, errors.UnknownParameter
}
