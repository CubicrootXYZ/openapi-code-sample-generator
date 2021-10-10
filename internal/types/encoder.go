package types

// Encoder defines an interface for encoding
type Encoder interface {
	EnocdeValue(ref string, value interface{}) (string, error)
	EnocdeParameter(name string, value interface{}) (string, error)
}

const (
	EncodingWwwUrlencode = "application/x-www-form-urlencoded"
	EncodingJSON         = "application/json"
	EncodingXML          = "application/xml"
)
