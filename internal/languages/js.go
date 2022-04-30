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
func (js *JS) GetSample(httpVerb string, path string, operation *openapi3.Operation, pathItem *openapi3.PathItem, document *openapi3.T) (*types.CodeSample, error) {
	js.usesToken = false
	codeInit := strings.Builder{} // js variable definitions
	codeExec := strings.Builder{} // js code

	parameters, err := js.extractor.GetParameters(operation.Parameters)
	if err != nil {
		return nil, err
	}

	secParameters, basicAuth, err := js.extractor.GetSecurity(operation, document)
	if err != nil {
		return nil, err
	}

	parameters.Query = append(parameters.Query, secParameters.Query...)
	parameters.Header = append(parameters.Header, secParameters.Header...)
	parameters.Path = append(parameters.Path, secParameters.Path...)
	parameters.Cookie = append(parameters.Cookie, secParameters.Cookie...)

	body, meta := js.getRequestBody(operation)
	body = js.filterToken(body)

	codeInit.WriteString("var request = new XMLHttpRequest();\n")

	// Set the url
	codeExec.WriteString("request.open(\"")
	codeExec.WriteString(httpVerb)
	codeExec.WriteString("\", \"https://")
	codeExec.WriteString(js.extractor.GetURL(operation, pathItem, document))
	codeExec.WriteString(js.extractor.GetPathExample(path, parameters.Path))
	if len(parameters.Query) > 0 {
		codeExec.WriteString("?")
		codeExec.WriteString(strings.TrimSuffix(js.filterToken(js.getQueryParams(parameters.Query)+"\""), " + \"\""))
	} else {
		codeExec.WriteString("\"")
	}
	codeExec.WriteString(", false);\n")

	if js.usesToken {
		codeInit.WriteString("var token = \"my secret token\"; // Put your token here\n")
	}

	// Set request headers
	for _, header := range parameters.Header {
		codeExec.WriteString("request.setRequestHeader(\"")
		codeExec.WriteString(header.Name)
		codeExec.WriteString("\", ")
		codeExec.WriteString(strings.TrimSuffix(strings.TrimPrefix("\""+js.filterToken(fmt.Sprint(header.Value)+"\""), "\"\" + "), " + \"\""))
		codeExec.WriteString(");\n")
	}

	if basicAuth {
		codeInit.WriteString("var username = \"my username\";\n")
		codeInit.WriteString("var password = \"*******\";\n")
		codeExec.WriteString("request.setRequestHeader(\"Authorization\", \"Basic \" + btoa(username + \":\" + password));\n")
	}

	if meta.Format != "" {
		codeExec.WriteString("request.setRequestHeader(\"Content-Type\", \"")
		codeExec.WriteString(meta.Format)
		if meta.FormData.OuterBoundary != nil {
			codeExec.WriteString(" boundary=")
			codeExec.WriteString(*meta.FormData.OuterBoundary)
		}
		codeExec.WriteString("\");\n")
	}

	// Set request body
	codeExec.WriteString("request.send(\"")
	codeExec.WriteString(body)
	codeExec.WriteString("\");\n")

	codeExec.WriteString("console.log(request.responseText);")
	codeInit.WriteString("\n")

	return &types.CodeSample{
		Lang:   types.LanguageJS,
		Label:  "JavaScript",
		Source: codeInit.String() + codeExec.String(),
	}, nil
}

func (js *JS) getQueryParams(params []*types.Parameter) string {
	query := strings.Builder{}
	encoder, ok := js.encoders[types.EncodingWwwUrlencode]

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

		query.WriteString(js.escape(encoded))
	}

	log.Debug(fmt.Sprintf("Wrote %d parameters to query", len(params)))

	return query.String()
}

func (js *JS) getRequestBody(operation *openapi3.Operation) (string, *types.FormattingMeta) {
	meta := &types.FormattingMeta{}
	if operation.RequestBody == nil || operation.RequestBody.Value == nil {
		return "", meta
	}

	value, format, err := js.extractor.GetRequestBody(operation.RequestBody.Value)
	if err != nil {
		log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
		return "", meta
	}

	meta.Format = format

	encoder, ok := js.encoders[strings.ToLower(format)]
	if !ok {
		log.Warn("Missing encoder for format: json")
		return "", meta
	}

	newValue, err := encoder.EnocdeValue(operation.RequestBody.Ref, value, meta)
	if err != nil {
		log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
		return "", meta
	}

	return js.escape(newValue), meta
}

func (p *JS) escape(text string) string {
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.ReplaceAll(text, "\r\n", "\\r\\n")

	return text
}

func (js *JS) filterToken(text string) string {

	if strings.Contains(text, "${TOKEN}") {
		js.usesToken = true
		text = strings.Replace(text, "${TOKEN}", "\" + token + \"", -1)
	}

	return text
}
