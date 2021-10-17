package errors

import "errors"

var ErrUnknownLanguage = errors.New("unknown language")

var ErrNoServer = errors.New("no server found")

var ErrUnknownSchema = errors.New("unknown schema")
var ErrUnknownParameter = errors.New("unknown parameter")
var ErrUnknownMediaType = errors.New("unknown media type")
var ErrUnsupportedDataType = errors.New("data type is not supported")

var ErrTypeAssertionFailed = errors.New("type assertion failed")
