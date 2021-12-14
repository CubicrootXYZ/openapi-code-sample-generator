package types

// Encoder defines an interface for encoding
type Encoder interface {
	EnocdeValue(ref string, value interface{}, meta *FormattingMeta) (string, error)
	EnocdeParameter(name string, value interface{}) (string, error)
}

const (
	EncodingWwwUrlencode = "application/x-www-form-urlencoded"
	EncodingJSON         = "application/json"
	EncodingJSONText     = "text/json"
	EncodingXML          = "application/xml"
	EncodingXMLText      = "text/xml"
	EncodingFormData     = "multipart/form-data"
	EncodingJSONPHP      = "application/json/php"
)
