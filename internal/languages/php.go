package languages

import (
	"fmt"
	"strings"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// Php holds information for php samples
type Php struct {
	encoders  map[string]types.Encoder
	extractor types.Extractor
}

// NewPhp returns a new php object
func NewPhp(encoders map[string]types.Encoder, extractor types.Extractor) types.Generator {
	return &Php{
		encoders:  encoders,
		extractor: extractor,
	}
}

// GetSample returns a php sample for the given operation
func (p *Php) GetSample(httpVerb string, path string, operation *openapi3.Operation, pathItem *openapi3.PathItem, document *openapi3.T) (*types.CodeSample, error) {
	codeInit := strings.Builder{} // php variable definitions
	codeExec := strings.Builder{} // php curl statement

	parameters, err := p.extractor.GetParameters(operation.Parameters)
	if err != nil {
		return nil, err
	}

	secParameters, basicAuth, err := p.extractor.GetSecurity(operation, document)
	if err != nil {
		return nil, err
	}

	parameters.Query = append(parameters.Query, secParameters.Query...)
	parameters.Header = append(parameters.Header, secParameters.Header...)
	parameters.Path = append(parameters.Path, secParameters.Path...)
	parameters.Cookie = append(parameters.Cookie, secParameters.Cookie...)

	body, meta := p.getRequestBody(operation)
	headers := p.getHeaderParams(parameters.Header, meta)
	cookies := p.getCookieParams(parameters.Cookie)

	// Basic curl stuff
	codeExec.WriteString("\n\n$curl = curl_init($url);\ncurl_setopt($curl, CURLOPT_RETURNTRANSFER, true);\n")

	// Set url
	codeInit.WriteString("<?php\n$url = \"")
	codeInit.WriteString(p.extractor.GetURL(operation, pathItem, document))
	codeInit.WriteString(p.extractor.GetPathExample(path, parameters.Path))
	if len(parameters.Query) > 0 {
		codeInit.WriteString("?")
		codeInit.WriteString(p.getQueryParams(parameters.Query))
	}
	codeInit.WriteString("\";\n")

	// Set request body
	if body != "" {
		// TODO use php arrays where possible and json_encode
		codeInit.WriteString("$data = ")
		codeInit.WriteString(body)
		codeInit.WriteString(";\n")

		codeExec.WriteString("curl_setopt($curl, CURLOPT_POSTFIELDS, $data);\n")
	}

	// Set request headers
	if headers != "array()" {
		codeInit.WriteString("$headers = ")
		codeInit.WriteString(headers)
		codeInit.WriteString(";\n")

		codeExec.WriteString("curl_setopt($curl, CURLOPT_HTTPHEADER, $headers);\n")
	}

	// Set request cookies
	if cookies != "" {
		codeInit.WriteString("$cookies = \"")
		codeInit.WriteString(cookies)
		codeInit.WriteString("\";\n")

		codeExec.WriteString("curl_setopt($curl, CURLOPT_COOKIE, $cookies);\n")
	}

	if basicAuth {
		codeInit.WriteString("$username = \"username\";\n")
		codeInit.WriteString("$password = \"password\";\n")

		codeExec.WriteString("curl_setopt($ch, CURLOPT_USERPWD, $username . \":\" . $password);\n")
	}

	codeExec.WriteString("curl_setopt($curl, CURLOPT_CUSTOMREQUEST, '")
	codeExec.WriteString(strings.ToUpper(httpVerb))
	codeExec.WriteString("');\n")

	codeExec.WriteString("var_dump(curl_exec($curl)); // Dumps the response\ncurl_close($curl);")

	return &types.CodeSample{
		Lang:   types.LanguagePhp,
		Label:  "PHP",
		Source: codeInit.String() + codeExec.String(),
	}, nil
}

func (p *Php) getQueryParams(params []*types.Parameter) string {
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

func (p *Php) getHeaderParams(params []*types.Parameter, meta *types.FormattingMeta) string {
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

	head.WriteString("\t\"")
	head.WriteString(p.writeFormatMeta(meta))
	head.WriteString("\",\n")

	head.WriteString(")")

	// TODO handle ${TOKEN}

	return head.String()
}

func (p *Php) getCookieParams(params []*types.Parameter) string {
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

func (p *Php) getRequestBody(operation *openapi3.Operation) (string, *types.FormattingMeta) {
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

	addQuotes := false
	encoder, ok := p.encoders[strings.ToLower(format)+"/php"]
	if !ok {
		addQuotes = true
		log.Debug("No PHP specific encoder, falling back to: " + format)
		encoder, ok = p.encoders[strings.ToLower(format)]
	}

	if ok {
		newValue, err := encoder.EnocdeValue(operation.RequestBody.Ref, value, meta)
		if err != nil {
			log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
			return "", meta
		}

		if !addQuotes {
			return newValue, meta
		}

		return "\"" + p.escape(newValue) + "\"", meta
	} else {
		log.Warn("Missing encoder for format: " + format)
	}

	return "", meta
}

func (p *Php) escape(text string) string {
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.ReplaceAll(text, "\r\n", "\\r\\n")

	return text
}

func (p *Php) writeFormatMeta(meta *types.FormattingMeta) string {
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
