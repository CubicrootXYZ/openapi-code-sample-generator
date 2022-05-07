package curl

import (
	_ "embed"
	"net/url"
	"strings"
	"text/template"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater"
)

//go:embed sample.tmpl
var templateString string

type Language struct {
}

func New() *Language {
	return &Language{}
}

func (Language *Language) Name() string {
	return "curl"
}

func (language *Language) GetTemplate() (*template.Template, error) {
	tmpl := template.New(language.Name())

	tmpl.Funcs(
		template.FuncMap{
			"escape":    escape,
			"urlencode": url.QueryEscape,
		},
	)

	tmpl, err := tmpl.Parse(templateString)
	if err != nil {
		return tmpl, err
	}

	return tmpl, nil
}

func (language *Language) GetAdditionals(data *templater.TemplateData) map[string]interface{} {
	return map[string]interface{}{}
}

func escape(value string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(value, `"`, `\"`, -1), "\n", "\\n", -1), "\r", "\\r", -1), "\t", "\\t", -1)
}
