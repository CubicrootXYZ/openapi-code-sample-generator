package languages

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// Curl holds information for curl samples
type Curl struct {
	encoders  map[string]types.Encoder
	extractor types.Extractor
}

// NewCurl returns a new curl object
func NewCurl(encoders map[string]types.Encoder, extractor types.Extractor) types.Generator {
	return &Curl{
		encoders:  encoders,
		extractor: extractor,
	}
}

// GetSample returns a curl sample for the given operation
func (c *Curl) GetSample(httpVerb string, path string, operation *openapi3.Operation, pathItem *openapi3.PathItem, document *openapi3.T) (*types.CodeSample, error) {
	cmd := strings.Builder{}
	parameters, err := c.extractor.GetParameters(operation.Parameters)
	if err != nil {
		return nil, err
	}

	secParameters, basicAuth, err := c.extractor.GetSecurity(operation, document)
	if err != nil {
		return nil, err
	}

	parameters.Query = append(parameters.Query, secParameters.Query...)
	parameters.Header = append(parameters.Header, secParameters.Header...)
	parameters.Path = append(parameters.Path, secParameters.Path...)
	parameters.Cookie = append(parameters.Cookie, secParameters.Cookie...)

	body, meta := c.getRequestBody(operation)

	cmd.WriteString("curl \"")
	cmd.WriteString(c.extractor.GetURL(operation, pathItem, document))
	cmd.WriteString(c.extractor.GetPathExample(path, parameters.Path))
	if len(parameters.Query) > 0 {
		cmd.WriteString("?")
		cmd.WriteString(c.getQueryParams(parameters.Query))
	}
	cmd.WriteString("\"")
	if basicAuth {
		cmd.WriteString(" -u username:password")
	}
	cmd.WriteString(c.getHeaderParams(parameters.Header))
	cmd.WriteString(c.getCookieParams(parameters.Cookie))
	cmd.WriteString(" -d \"")
	cmd.WriteString(body)
	cmd.WriteString("\" ")
	if meta != nil {
		cmd.WriteString(c.writeFormatMeta(meta))
	}
	cmd.WriteString("-X ")
	cmd.WriteString(strings.ToUpper(httpVerb))

	return &types.CodeSample{
		Lang:   types.LanguageCurl,
		Label:  "curl",
		Source: cmd.String(),
	}, nil
}

func (c *Curl) getQueryParams(params []*types.Parameter) string {
	query := strings.Builder{}
	encoder, ok := c.encoders[types.EncodingWwwUrlencode]

	if !ok {
		log.Warn("Missing encoder for format: " + types.EncodingWwwUrlencode)
		return ""
	}

	for i, param := range params {
		if param == nil {
			continue
		}

		encoded, err := encoder.EnocdeParameter(param.Name, param.Value)
		if err != nil {
			continue
		}

		if i != 0 {
			query.WriteString("&")
		}

		query.WriteString(c.escape(encoded))
	}

	log.Debug(fmt.Sprintf("Wrote %d parameters to query", len(params)))

	return query.String()
}

func (c *Curl) getHeaderParams(params []*types.Parameter) string {
	head := strings.Builder{}
	encoder, ok := c.encoders[types.EncodingWwwUrlencode]

	if !ok {
		log.Warn("Missing encoder for format: " + types.EncodingWwwUrlencode)
		return ""
	}

	for _, param := range params {
		if param == nil {
			continue
		}

		value, err := encoder.EnocdeValue("", param.Value, nil)
		if err != nil {
			log.Info(fmt.Sprintf("Skipped header parameter %s due to: %s", param.Name, err.Error()))
			continue
		}

		head.WriteString(" -H \"")
		head.WriteString(c.escape(param.Name))
		head.WriteString(": ")
		head.WriteString(c.escape(value))
		head.WriteString("\"")
	}

	return head.String()
}

func (c *Curl) getCookieParams(params []*types.Parameter) string {
	head := strings.Builder{}
	encoder, ok := c.encoders[types.EncodingWwwUrlencode]

	if !ok {
		log.Warn("Missing encoder for format: " + types.EncodingWwwUrlencode)
		return ""
	}

	for _, param := range params {
		if param == nil {
			continue
		}

		value, err := encoder.EnocdeParameter(param.Name, param.Value)
		if err != nil {
			log.Info(fmt.Sprintf("Skipped cookie parameter %s due to: %s", param.Name, err.Error()))
			continue
		}

		head.WriteString(" -b \"")
		head.WriteString(c.escape(value))
		head.WriteString("\"")
	}

	return head.String()
}

func (c *Curl) getRequestBody(operation *openapi3.Operation) (string, *types.FormattingMeta) {
	meta := &types.FormattingMeta{}
	if operation.RequestBody == nil || operation.RequestBody.Value == nil {
		return "", meta
	}

	value, format, err := c.extractor.GetRequestBody(operation.RequestBody.Value)
	if err != nil {
		log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
		return "", meta
	}

	meta.Format = format

	if encoder, ok := c.encoders[strings.ToLower(format)]; ok {
		newValue, err := encoder.EnocdeValue(operation.RequestBody.Ref, value, meta)
		if err != nil {
			log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
			return "", meta
		}
		return c.escape(newValue), meta
	} else {
		log.Warn("Missing encoder for format: " + format)
	}

	return "", meta
}

func (c *Curl) escape(text string) string {
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.ReplaceAll(text, "\r\n", "\\r\\n")
	return text
}

func (c *Curl) writeFormatMeta(meta *types.FormattingMeta) string {
	if meta == nil {
		return ""
	}

	cmd := strings.Builder{}

	if meta.Format != "" {
		cmd.WriteString("-H \"Content-Type: ")
		cmd.WriteString(meta.Format)

		if meta.FormData.OuterBoundary != nil {
			cmd.WriteString(" boundary=")
			cmd.WriteString(*meta.FormData.OuterBoundary)
		}
		cmd.WriteString("\" ")
	}

	return cmd.String()
}
