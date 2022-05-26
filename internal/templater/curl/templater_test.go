package curl

import (
	"testing"

	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/encoding"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/extractor"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/log"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater"
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.Verbose = true
	m.Run()
}

func testEndpoint1() *templater.Endpoint {
	endpoint := templater.Endpoint{}
	endpoint.HTTPVerb = "POST"
	endpoint.Path = "/random_#+!§$%&/()=/path"
	endpoint.OpenAPI.Operation = &openapi3.Operation{
		Parameters: openapi3.Parameters{&openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name: "param1",
				In:   "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
				Required: true,
			},
		}, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name: "param2",
				In:   "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
				Required: false,
			},
		}, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name: "param3",
				In:   "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   "string",
						Format: "date-time",
					},
				},
				Required: true,
			},
		}, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name: "param4",
				In:   "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "boolean",
					},
				},
				Required: true,
			},
		}, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name: "param5",
				In:   "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "array",
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
					},
				},
				Required: true,
			},
		}, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name: "param6",
				In:   "query",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "object",
						AnyOf: openapi3.SchemaRefs{
							nil,
							&openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: "object",
									Properties: openapi3.Schemas{"param6-sub": &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: "array",
											Items: &openapi3.SchemaRef{
												Value: &openapi3.Schema{
													Type: "string",
												},
											},
										},
									},
									},
								},
							},
						},
					},
				},
				Required: true,
			},
		}},
	}
	endpoint.OpenAPI.PathItem = &openapi3.PathItem{
		Post: endpoint.OpenAPI.Operation,
	}
	endpoint.OpenAPI.Document = &openapi3.T{
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL: "example.com",
			},
		},
	}

	return &endpoint
}

func TestTemplater_Template_Curl(t *testing.T) {
	templater := templater.NewTemplater(encoding.Encoders(),
		extractor.NewOpenAPIExtractor(),
		map[types.Language]templater.Language{types.LanguageCurl: New()})

	result, err := templater.Template(types.LanguageCurl, testEndpoint1())

	require.NoError(t, err)
	assert.Equal(t, types.LanguageCurl, result.Lang)
	assert.Equal(t, New().Name(), result.Label)
	assert.Equal(t, `curl https://example.com/random_#+!§$%&/()=/path?param1=1234&param3=2022-01-01T15%3A00%3A14Z&param4=false&param5%5B%5D=example-string&param6%5Bparam6-sub%5D%5B%5D=example-string -X POST`, result.Source)
}
