package encoding

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/errors"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

// FormDataEncode groups multipart/form-data features
type FormDataEncode struct {
}

// EnocdeParameter encodes the given parameter and its value to multipart/form-data
func (f *FormDataEncode) EnocdeParameter(name string, value interface{}) (string, error) {
	return f.EnocdeValue("", map[string]interface{}{name: value}, nil)
}

// EnocdeValue encodes a single value to amultipart/form-data
func (f *FormDataEncode) EnocdeValue(ref string, value interface{}, meta *types.FormattingMeta) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if newValue, ok := value.(map[interface{}]interface{}); ok {
		for key, value := range newValue {
			f.writePart(fmt.Sprint(key), value, writer)
		}
		writer.Close()

		return f.prefix(writer.Boundary(), meta) + body.String(), nil
	}

	if newValue, ok := value.(map[string]interface{}); ok {
		for key, value := range newValue {
			f.writePart(key, value, writer)
		}
		writer.Close()

		return f.prefix(writer.Boundary(), meta) + body.String(), nil
	}

	log.Error(fmt.Sprint(value))

	log.Warn(fmt.Sprintf("Data type %T is not supported for multipart/form-data", value))
	return "", errors.ErrUnsupportedDataType
}

func (f *FormDataEncode) prefix(boundary string, meta *types.FormattingMeta) string {
	if meta != nil {
		meta.FormData.OuterBoundary = &boundary
	}
	return "Content-Type: multipart/form-data; boundary=" + boundary + "\r\n\r\n"
}

func (f *FormDataEncode) writePart(key string, value interface{}, writer *multipart.Writer) {
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "text/plain")
	metadataHeader.Set("Content-ID", key)
	part, _ := writer.CreatePart(metadataHeader)

	_, ok1 := value.(map[interface{}]interface{})
	_, ok2 := value.(map[string]interface{})

	if ok1 || ok2 {
		newVal, err := f.EnocdeValue("", value, nil)
		if err == nil {
			_, _ = part.Write([]byte(newVal))
		} else {
			log.Warn(err.Error())
		}
	} else {
		_, _ = part.Write([]byte(fmt.Sprint(value)))
	}
}
