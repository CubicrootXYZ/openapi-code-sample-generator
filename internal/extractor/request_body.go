package extractor

import (
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/helper"
	"openapi-code-sample-generator/internal/log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetRequestBody returns all parameters inside the request body
func (o *openAPIExtractor) GetRequestBody(body *openapi3.RequestBody) (value interface{}, format string, err error) {
	for encoding, mediaType := range body.Content {
		if mediaType == nil {
			continue
		}
		value, err = o.getMediaTypeValue(mediaType, strings.ToLower(encoding))
		if err != nil {
			continue
		}
		format = strings.ToLower(encoding)
		break
	}

	return
}

func (o *openAPIExtractor) getMediaTypeValue(mediaType *openapi3.MediaType, format string) (interface{}, error) {
	if !helper.IsNil(mediaType.Example) {
		log.Debug("using param example value")
		return mediaType.Example, nil
	}

	if mediaType.Schema != nil && mediaType.Schema.Value != nil {
		val, err := o.GetExampleValueForSchema(mediaType.Schema.Value, format)
		if err == nil {
			return val, nil
		}
	}

	return nil, errors.UnknownMediaType
}