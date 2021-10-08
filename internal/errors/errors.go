package errors

import "errors"

var UnknownLanguage = errors.New("unknown language")

var NoServer = errors.New("no server found")

var UnknownSchema = errors.New("unknown schema")
var UnknownParameter = errors.New("unknown parameter")
var UnknownMediaType = errors.New("unknown media type")

var TypeAssertionFailed = errors.New("type assertion failed")
