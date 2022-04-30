package languages

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// JS holds information for js samples
type JS struct {
	encoders  map[string]types.Encoder
	extractor types.Extractor
	usesToken bool
}

// NewJS returns a new JS object
func NewJS(encoders map[string]types.Encoder, extractor types.Extractor) types.Generator {
	return &JS{
		encoders:  encoders,
		extractor: extractor,
	}
}

// GetSample returns a js sample for the given operation
func (p *JS) GetSample(httpVerb string, path string, operation *openapi3.Operation, pathItem *openapi3.PathItem, document *openapi3.T) (*types.CodeSample, error) {
	p.usesToken = false
	codeInit := strings.Builder{} // js variable definitions
	codeExec := strings.Builder{} // js curl statement

	parameters, err := p.extractor.GetParameters(operation.Parameters)
	if err != nil {
		return nil, err
	}

	secParameters, basicAuth, err := p.extractor.GetSecurity(operation, document)
	if err != nil {
		return nil, err
	}

	parameters.Query = append(parameters.Query, parameters.Query...)
	parameters.Header = append(parameters.Header, secParameters.Header...)
	parameters.Path = append(parameters.Path, secParameters.Path...)
	parameters.Cookie = append(parameters.Cookie, secParameters.Cookie...)

	body, _ := p.getRequestBody(operation)
	body = p.filterToken(body)

	// Set the url
	codeExec.WriteString("request.open(\"")
	codeExec.WriteString(httpVerb)
	codeExec.WriteString("\", \"")
	codeExec.WriteString(p.extractor.GetURL(operation, pathItem, document))
	codeExec.WriteString(p.extractor.GetPathExample(path, parameters.Path))
	if len(parameters.Query) > 0 {
		codeInit.WriteString("?")
		codeInit.WriteString(p.getQueryParams(parameters.Query))
	}
	codeExec.WriteString(", false);\"\n")

	if p.usesToken {
		codeInit.WriteString("var token = \"my secret token\"; // Put your token here\n")
	}

	// Set request headers
	for _, header := range parameters.Header {
		codeExec.WriteString("request.setRequestHeader(\"")
		codeExec.WriteString(header.Name)
		codeExec.WriteString("\", ")
		codeExec.WriteString(strings.TrimSuffix(strings.TrimPrefix("\""+p.filterToken(fmt.Sprint(header.Value)+"\""), "\"\" + "), " + \"\""))
		codeExec.WriteString(");\n")
	}

	if basicAuth {
		codeInit.WriteString("var username = \"my username\";\n")
		codeInit.WriteString("var password = \"*******\";\n")
		codeExec.WriteString("request.setRequestHeader(\"Authorization\", \"Basic \" + btoa(username + \":\" + password));\n")
	}

	// Set request body
	codeExec.WriteString("request.send(\"")
	codeExec.WriteString(body)
	codeExec.WriteString("\");\n")

	codeExec.WriteString("console.log(request.responseText)")

	return &types.CodeSample{
		Lang:   types.LanguageJS,
		Label:  "JavaScript",
		Source: codeInit.String() + codeExec.String(),
	}, nil
}

func (p *JS) getQueryParams(params []*types.Parameter) string {
	query := strings.Builder{}
	encoder, ok := p.encoders[types.EncodingWwwUrlencode]

	if !ok {
		log.Warn("Missing encoder for format: " + types.EncodingWwwUrlencode)
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

		query.WriteString(p.escape(encoded))
	}

	log.Debug(fmt.Sprintf("Wrote %d parameters to query", len(params)))

	return query.String()
}

func (p *JS) getHeaderParams(params []*types.Parameter, meta *types.FormattingMeta) string {
	head := strings.Builder{}
	encoder, ok := p.encoders[types.EncodingWwwUrlencode]

	if !ok {
		log.Warn("Missing encoder for format: " + types.EncodingWwwUrlencode)
		return "array()"
	}

	head.WriteString("array(\n")

	for _, param := range params {
		if param == nil {
			continue
		}

		value, err := encoder.EnocdeValue("", param.Value, nil)
		if err != nil {
			log.Info(fmt.Sprintf("Skipped header parameter %s due to: %s", param.Name, err.Error()))
			continue
		}

		head.WriteString("\t\"")
		head.WriteString(p.escape(param.Name))
		head.WriteString(": ")
		head.WriteString(p.escape(value))
		head.WriteString("\",\n")
	}

	encoding := p.writeFormatMeta(meta)

	if encoding != "" {
		head.WriteString("\t\"")
		head.WriteString(encoding)
		head.WriteString("\",\n")
	}

	head.WriteString(")")

	// TODO handle ${TOKEN} also in cookie and params

	return head.String()
}

func (p *JS) getCookieParams(params []*types.Parameter) string {
	cookie := strings.Builder{}
	encoder, ok := p.encoders[types.EncodingWwwUrlencode]

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

		cookie.WriteString(p.escape(value))
		cookie.WriteString(";")
	}

	return cookie.String()
}

func (p *JS) getRequestBody(operation *openapi3.Operation) (string, *types.FormattingMeta) {
	meta := &types.FormattingMeta{}
	if operation.RequestBody == nil || operation.RequestBody.Value == nil {
		return "", meta
	}

	value, format, err := p.extractor.GetRequestBody(operation.RequestBody.Value)
	if err != nil {
		log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
		return "", meta
	}

	meta.Format = format

	encoder, ok := p.encoders[types.EncodingJSON]
	if !ok {
		log.Warn("Missing encoder for format: json")
		return "", meta
	}

	jsEncodedValue, err := encoder.EnocdeValue(operation.RequestBody.Ref, value, meta)
	if err == nil {
		switch strings.ToLower(format) {
		case types.EncodingJSON, types.EncodingJSONText:
			return "json_encode(" + jsEncodedValue + ")", meta
		}
	}
	log.Warn("Failed js encoding value, fallbacking. Error was: " + err.Error())

	newValue, err := encoder.EnocdeValue(operation.RequestBody.Ref, value, meta)
	if err != nil {
		log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
		return "", meta
	}

	return "\"" + p.escape(newValue) + "\"", meta
}

func (p *JS) escape(text string) string {
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.ReplaceAll(text, "\r\n", "\\r\\n")

	return text
}

func (p *JS) writeFormatMeta(meta *types.FormattingMeta) string {
	if meta == nil {
		return ""
	}

	code := strings.Builder{}

	if meta.Format != "" {
		code.WriteString("Content-Type: ")
		code.WriteString(meta.Format)

		if meta.FormData.OuterBoundary != nil {
			code.WriteString(" boundary=")
			code.WriteString(*meta.FormData.OuterBoundary)
		}
	}

	return code.String()
}

func (p *JS) filterToken(text string) string {

	if strings.Contains(text, "${TOKEN}") {
		p.usesToken = true
		text = strings.Replace(text, "${TOKEN}", "\" + token + \"", -1)
	}

	return text
}
