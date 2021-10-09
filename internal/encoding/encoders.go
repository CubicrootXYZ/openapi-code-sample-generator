package encoding

import "openapi-code-sample-generator/internal/types"

// Encoders returns all available encoders
func Encoders() map[string]types.Encoder {
	encoders := make(map[string]types.Encoder)
	encoders[types.EncodingWwwUrlencode] = &URLEncode{}
	encoders[types.EncodingJSON] = &JSONEncode{}

	return encoders
}
