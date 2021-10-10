package encoding

import "github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

// Encoders returns all available encoders
func Encoders() map[string]types.Encoder {
	encoders := make(map[string]types.Encoder)
	encoders[types.EncodingWwwUrlencode] = &URLEncode{}
	encoders[types.EncodingJSON] = &JSONEncode{}
	encoders[types.EncodingXML] = &XMLEncode{}

	return encoders
}
