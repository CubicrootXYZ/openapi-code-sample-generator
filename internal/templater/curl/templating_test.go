package curl

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

func TestTemplater_Template_Curl(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New()})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint1())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New().Name(), result.Label)
	assert.Equal(t, `curl https://example.com/random_#+!§$%&/()=/path?param1=1234&param3=2022-01-01T15%3A00%3A14Z&param4=false&param5%5B%5D=example-string&param6%5Bparam6-sub%5D%5B%5D=example-string -X POST`, result.Source)
}

func TestTemplater_Template_CurlJsonBody(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New()})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint2())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New().Name(), result.Label)
	assert.Equal(t, `curl https://example.com/random_#+!§$%&/()=/path -X POST -H "Content-Type: application/json " -d "{\"param6-sub\":[\"example-string\"]}"`, result.Source)
}

func TestTemplater_Template_CurlXmlBody(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New()})

	result, err := templater.Template(types.LanguageCurl, test.Endpoint3())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New().Name(), result.Label)
	assert.Equal(t, `curl https://example.com/random_#+!§$%&/()=/path -X POST -H "Content-Type: application/xml " -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?><doc><param6-sub>example-string</param6-sub></doc>"`, result.Source)
}
