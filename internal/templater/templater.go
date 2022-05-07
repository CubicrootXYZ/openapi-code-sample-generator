package templater

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
	"github.com/getkin/kin-openapi/openapi3"
)

var ErrMissingEncoder = errors.New("missing encoder")
var ErrUnknownLanguage = errors.New("language is not known")

// Templater templates a template
type Templater interface {
	Template(language types.Language, endpoint *Endpoint) (string, error)
}

// TemplateData holds all data available in templating
type TemplateData struct {
	Parameters         *types.Parameters      // Different parameter types
	SecurityParameters *types.Parameters      // Security related parameters (authentication), their value is either "${TOKEN}" or "Bearer ${TOKEN}"
	Path               string                 // Example path
	QueryParamsString  string                 // Example query params as urlencoded string
	URL                string                 // Example URL
	HTTPVerb           string                 // HTTP Method
	Body               interface{}            // The body as golang data type
	BodyString         string                 // The body formatted in the endpoints "accept" format	Meta               *types.FormattingMeta
	BasicAuth          bool                   // True if basic Auth is enabled
	Additionals        map[string]interface{} // Language specific data
	Formatting         *types.FormattingMeta  // Metadata about the format used
}

// Language defines an interface for each language to be used with the templater
type Language interface {
	GetAdditionals(data *TemplateData) map[string]interface{} // Allows the language to set additional data to the template data
	GetTemplate() (*template.Template, error)                 // Must return the template for the code sample
	// TODO add custom escape method
}

// Endpoint defines a openapi operation and meta information for templating
type Endpoint struct {
	HTTPVerb string
	Path     string
	OpenAPI  struct {
		Operation *openapi3.Operation
		PathItem  *openapi3.PathItem
		Document  *openapi3.T
	}
}

func NewEndpoint(HTTPVerb string, path string, operation *openapi3.Operation, pathItem *openapi3.PathItem, document *openapi3.T) *Endpoint {
	return &Endpoint{
		HTTPVerb: HTTPVerb,
		Path:     path,
		OpenAPI: struct {
			Operation *openapi3.Operation
			PathItem  *openapi3.PathItem
			Document  *openapi3.T
		}{
			Operation: operation,
			PathItem:  pathItem,
			Document:  document,
		},
	}
}

type templater struct {
	encoders  map[string]types.Encoder
	extractor types.Extractor
	languages map[types.Language]Language
}

// NewTemplater creates a new templating instance
func NewTemplater(encoders map[string]types.Encoder, extractor types.Extractor, languages map[types.Language]Language) Templater {
	return &templater{
		encoders:  encoders,
		extractor: extractor,
		languages: languages,
	}
}

// Template actually templates the template
func (template *templater) Template(lang types.Language, endpoint *Endpoint) (string, error) {
	language, ok := template.languages[lang]
	if !ok {
		return "", ErrUnknownLanguage
	}

	templateData, err := template.extractTemplateData(endpoint)
	if err != nil {
		return "", err
	}

	templateData.Additionals = language.GetAdditionals(templateData)
	tmpl, err := language.GetTemplate()
	if err != nil {
		return "", err
	}

	buffer := bytes.Buffer{}
	err = tmpl.Execute(&buffer, &templateData)
	if err != nil {
		return "", err
	}

	return buffer.String(), err
}

func (template *templater) extractTemplateData(endpoint *Endpoint) (*TemplateData, error) {
	templateData := TemplateData{}
	params, err := template.extractor.GetParameters(endpoint.OpenAPI.Operation.Parameters)
	if err != nil {
		return nil, err
	}

	secParameters, basicAuth, err := template.extractor.GetSecurity(endpoint.OpenAPI.Operation, endpoint.OpenAPI.Document)
	if err != nil {
		return nil, err
	}

	requestBodyInterface, requestBodyString, meta, err := template.extractRequestBody(endpoint.OpenAPI.Operation)
	if err != nil {
		return nil, err
	}

	url := template.extractor.GetURL(endpoint.OpenAPI.Operation, endpoint.OpenAPI.PathItem, endpoint.OpenAPI.Document)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	templateData.SecurityParameters = &secParameters
	templateData.HTTPVerb = endpoint.HTTPVerb
	templateData.Parameters = &params
	templateData.BasicAuth = basicAuth
	templateData.Formatting = meta
	templateData.Body = requestBodyInterface
	templateData.BodyString = requestBodyString
	templateData.Path = template.extractor.GetPathExample(endpoint.Path, params.Path)
	templateData.URL = url
	templateData.QueryParamsString = template.buildExampleQueryParams(params.Query)
	// TODO scan for token

	return &templateData, nil
}

func (template *templater) extractRequestBody(operation *openapi3.Operation) (interface{}, string, *types.FormattingMeta, error) {
	meta := &types.FormattingMeta{}
	if operation.RequestBody == nil || operation.RequestBody.Value == nil {
		return nil, "", meta, nil
	}

	value, format, err := template.extractor.GetRequestBody(operation.RequestBody.Value)
	if err != nil {
		log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
		return nil, "", meta, err
	}

	meta.Format = format

	if encoder, ok := template.encoders[strings.ToLower(format)]; ok {
		newValue, err := encoder.EnocdeValue(operation.RequestBody.Ref, value, meta)
		if err != nil {
			log.Warn(fmt.Sprintf("Request body parsing failed: %s", err.Error()))
			return nil, "", meta, err
		}
		return value, newValue, meta, nil
	}

	log.Warn("Missing encoder for format: " + format)
	return nil, "", meta, ErrMissingEncoder
}

func (template *templater) buildExampleQueryParams(params []*types.Parameter) string {
	query := strings.Builder{}
	encoder, ok := template.encoders[types.EncodingWwwUrlencode]

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

		query.WriteString(encoded)
	}

	log.Debug(fmt.Sprintf("Wrote %d parameters to query", len(params)))

	return query.String()
}
