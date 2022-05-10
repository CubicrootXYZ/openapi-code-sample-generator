package php

import (
	_ "embed"
	"errors"
	"net/url"
	"strings"
	"text/template"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
)

//go:embed sample.gtpl
var templateString string

var ErrFormatNotSupported = errors.New("format not supported")

type Language struct {
	phpEncoder types.Encoder
}

func New(phpEncoder types.Encoder) *Language {
	return &Language{
		phpEncoder: phpEncoder,
	}
}

func (Language *Language) Name() string {
	return "PHP"
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

	phpRequestBody, err := language.getRequestBody(data)
	if err == nil {
		additionals["customRequestBody"] = phpRequestBody
	} else {
		log.Debug("Custom PHP body not supported for " + data.Formatting.Format)
	}

	return additionals
}

func (language *Language) getRequestBody(data *templater.TemplateData) (string, error) {
	formatSupported := false
	for _, formatDef := range []string{types.EncodingJSON, types.EncodingJSONText} {
		if strings.ToLower(data.Formatting.Format) == formatDef {
			formatSupported = true
		}
	}

	if !formatSupported {
		return "", ErrFormatNotSupported
	}

	value, err := language.phpEncoder.EnocdeValue("", data.Body, data.Formatting)
	if err == nil {
		switch strings.ToLower(data.Formatting.Format) {
		case types.EncodingJSON, types.EncodingJSONText:
			return "json_encode(" + value + ")", nil
		}
	}
	log.Warn("Failed php encoding value, fallbacking. Error was: " + err.Error())

	return "", ErrFormatNotSupported
}

func escape(value string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(value, `"`, `\"`, -1), "\n", "\\n", -1), "\r", "\\r", -1), "\t", "\\t", -1)
}

func escapeQuotes(value string) string {
	return strings.Replace(value, `"`, `\"`, -1)
}

func tokenStringToPHP(token string) string {
	if token == "${TOKEN}" {
		return "$token"
	} else if strings.HasSuffix(token, "${TOKEN}") {
		return "\"" + strings.TrimSuffix(token, "${TOKEN}") + "\" . $token"
	} else if strings.HasPrefix(token, "${TOKEN}") {
		return "$token . \"" + strings.TrimPrefix(token, "${TOKEN}") + "\""
	} else {
		return "$token"
	}
}
