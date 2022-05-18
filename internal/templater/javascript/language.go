package php

import (
	_ "embed"
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

	return additionals
}

func escape(value string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(value, `"`, `\"`, -1), "\n", "\\n", -1), "\r", "\\r", -1), "\t", "\\t", -1)
}

func escapeQuotes(value string) string {
	return strings.Replace(value, `"`, `\"`, -1)
}

func tokenStringToPHP(token string) string {
	if token == "${TOKEN}" {
		return "token"
	} else if strings.HasSuffix(token, "${TOKEN}") {
		return "\"" + escape(strings.TrimSuffix(token, "${TOKEN}")) + "\" + $token"
	} else if strings.HasPrefix(token, "${TOKEN}") {
		return "token + \"" + escape(strings.TrimPrefix(token, "${TOKEN}")) + "\""
	} else {
		return "token"
	}
}
