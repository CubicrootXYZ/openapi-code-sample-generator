package helper

import (
	"openapi-code-sample-generator/internal/errors"
	"openapi-code-sample-generator/internal/log"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetRequestBody returns all parameters inside the request body
func GetRequestBody(body *openapi3.RequestBody) (value interface{}, format string, err error) {
	for encoding, mediaType := range body.Content {
		log.Error(encoding)
		if mediaType == nil {
			continue
		}
		value, err = getMediaTypeValue(mediaType)
		if err != nil {
			continue
		}
		format = encoding
		break
	}

	return
}

func getMediaTypeValue(mediaType *openapi3.MediaType) (interface{}, error) {
	if !IsNil(mediaType.Example) {
		log.Debug("using param example value")
		return mediaType.Example, nil
	}

	if mediaType.Schema != nil && mediaType.Schema.Value != nil {
		val, err := GetExampleValueForSchema(mediaType.Schema.Value)
		if err == nil {
			return val, nil
		}
	}

	return nil, errors.UnknownMediaType
}
