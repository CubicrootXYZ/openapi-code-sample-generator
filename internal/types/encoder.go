package types

// Encoder defines an interface for encoding
type Encoder interface {
	EnocdeValue(value interface{}) (string, error)
	EnocdeParameter(name string, value interface{}) (string, error)
}

const (
	WwwUrlencode = "application/x-www-form-urlencoded"
)
