package php

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

func TestTemplater_Template_Php(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New(&encoding.PhpEncode{})})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint1())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New(&encoding.PhpEncode{}).Name(), result.Label)
	assert.Equal(t, "$url = \"https://example.com/random_#+!§$%&/()=/path?param1=1234&param3=2022-01-01T15%3A00%3A14Z&param4=false&param5%5B%5D=example-string&param6%5Bparam6-sub%5D%5B%5D=example-string\";\n\n$curl = curl_init($url);\ncurl_setopt($curl, CURLOPT_CUSTOMREQUEST, \"POST\");\ncurl_setopt($curl, CURLOPT_RETURNTRANSFER, true);\nvar_dump(curl_exec($curl)); // Dumps the response\ncurl_close($curl);", result.Source)
}

func TestTemplater_Template_PhpJsonBody(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New(&encoding.PhpEncode{})})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint2())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New(&encoding.PhpEncode{}).Name(), result.Label)
	assert.Equal(t, "$url = \"https://example.com/random_#+!§$%&/()=/path\";\n$data = json_encode(array(\n\t\"param6-sub\" => array(\n\t\t\"example-string\",\n\t),\n));\n\n\n$curl = curl_init($url);\ncurl_setopt($curl, CURLOPT_CUSTOMREQUEST, \"POST\");\ncurl_setopt($curl, CURLOPT_RETURNTRANSFER, true);\ncurl_setopt($curl, CURLOPT_POSTFIELDS, $data);\nvar_dump(curl_exec($curl)); // Dumps the response\ncurl_close($curl);", result.Source)
}

func TestTemplater_Template_PhpXmlBody(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New(&encoding.PhpEncode{})})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint3())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New(&encoding.PhpEncode{}).Name(), result.Label)
	assert.Equal(t, "$url = \"https://example.com/random_#+!§$%&/()=/path\";\n$data = \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?><doc><param6-sub>example-string</param6-sub></doc>\";\n\n$curl = curl_init($url);\ncurl_setopt($curl, CURLOPT_CUSTOMREQUEST, \"POST\");\ncurl_setopt($curl, CURLOPT_RETURNTRANSFER, true);\ncurl_setopt($curl, CURLOPT_POSTFIELDS, $data);\nvar_dump(curl_exec($curl)); // Dumps the response\ncurl_close($curl);", result.Source)
}
