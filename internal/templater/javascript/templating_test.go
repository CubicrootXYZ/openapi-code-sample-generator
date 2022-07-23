package javascript

import (
	"testing"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/encoding"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/extractor"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater"
	test "github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater/_test"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.Verbose = true
	m.Run()
}

func TestTemplater_Template_JS(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New()})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint1())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New().Name(), result.Label)
	assert.Equal(t, "\nvar url = \"https://example.com/random_#+!§$%&/()=/path?param1=1234&param3=2022-01-01T15%3A00%3A14Z&param4=false&param5%5B%5D=example-string&param6%5Bparam6-sub%5D%5B%5D=example-string\";\n\nvar request = new XMLHttpRequest();\nrequest.open(\"POST\", url);\n\nrequest.send(\"\");\nconsole.log(request.responseText);", result.Source)
}

func TestTemplater_Template_JSJsonBody(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New()})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint2())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New().Name(), result.Label)
	assert.Equal(t, "\nvar url = \"https://example.com/random_#+!§$%&/()=/path\";\n\nvar request = new XMLHttpRequest();\nrequest.open(\"POST\", url);\nrequest.setRequestHeader(\"Content-Type\", \"application/json\");\n\nrequest.send({\n\t\"param6-sub\": [\n\t\t\"example-string\"\n\t]\n});\nconsole.log(request.responseText);", result.Source)
}

func TestTemplater_Template_JSXmlBody(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New()})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint3())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New().Name(), result.Label)
	assert.Equal(t, "\nvar url = \"https://example.com/random_#+!§$%&/()=/path\";\n\nvar request = new XMLHttpRequest();\nrequest.open(\"POST\", url);\nrequest.setRequestHeader(\"Content-Type\", \"application/xml\");\n\nrequest.send(\"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?><doc><param6-sub>example-string</param6-sub></doc>\");\nconsole.log(request.responseText);", result.Source)
}
