package javascript

import (
	_ "embed"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"text/template"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater"
)

//go:embed sample.gtpl
var templateString string

var ErrFormatNotSupported = errors.New("format not supported")

type Language struct {
}

func New() *Language {
	return &Language{}
}

func (Language *Language) Name() string {
	return "JavaScript"
}

func (language *Language) GetTemplate() (*template.Template, error) {
	tmpl := template.New(language.Name())

	tmpl.Funcs(
		template.FuncMap{
			"escape":       escape,
			"urlencode":    url.QueryEscape,
			"converttoken": tokenStringToPHP,
			"escapeQuotes": escapeQuotes,
		},
	)

	tmpl, err := tmpl.Parse(templateString)
	if err != nil {
		return tmpl, err
	}

	return tmpl, nil
}

func (language *Language) GetAdditionals(data *templater.TemplateData) map[string]interface{} {
	additionals := make(map[string]interface{})

	if data.Formatting.Format == "application/json" || data.Formatting.Format == "text/json" {
		jsonData, err := json.MarshalIndent(data.Body, "", "\t")
		if err == nil {
			additionals["jsBody"] = string(jsonData)
		}
	}

	return additionals
}

func escape(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, `"`, `\"`), "\n", "\\n"), "\r", "\\r"), "\t", "\\t")
}

func escapeQuotes(value string) string {
	return strings.ReplaceAll(value, `"`, `\"`)
}

func tokenStringToPHP(token string) string {
	if token == "${TOKEN}" {
		return "token"
	} else if strings.HasSuffix(token, "${TOKEN}") {
		return "\"" + escape(strings.TrimSuffix(token, "${TOKEN}")) + "\" + token"
	} else if strings.HasPrefix(token, "${TOKEN}") {
		return "token + \"" + escape(strings.TrimPrefix(token, "${TOKEN}")) + "\""
	} else {
		return "token"
	}
}
