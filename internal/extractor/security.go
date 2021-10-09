package extractor

import (
	"openapi-code-sample-generator/internal/types"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetSecurity returns security related parameters
func (o *openAPIExtractor) GetSecurity(operation *openapi3.Operation, document *openapi3.T) (parameters types.Parameters, basicAuth bool, err error) {
	params := types.Parameters{}
	params.Query = make([]*types.Parameter, 0)
	params.Header = make([]*types.Parameter, 0)
	params.Path = make([]*types.Parameter, 0)
	params.Cookie = make([]*types.Parameter, 0)
	basicAuth = false

	requirements := operation.Security

	if requirements != nil {
		for _, requirement := range *requirements {
			for name := range requirement {
				security := o.getSecuritySchema(name, document)
				if security == nil {
					continue
				}
				switch strings.ToLower(security.Type) {
				case "http":
					switch strings.ToLower(security.Scheme) {
					case "basic":
						basicAuth = true
					case "bearer":
						params.Header = append(params.Header, &types.Parameter{
							Name:  "Authorization",
							Value: "Bearer ${TOKEN}",
						})
					}
				case "apikey":
					switch strings.ToLower(security.In) {
					case "query":
						params.Query = append(params.Query, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					case "cookie":
						params.Cookie = append(params.Cookie, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					case "header":
						params.Header = append(params.Header, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					}
				case "openidconnect", "oauth2":
					params.Header = append(params.Header, &types.Parameter{
						Name:  "Authorization",
						Value: "Bearer ${TOKEN}",
					})
				}
			}
		}
	}

	return
}

func (o *openAPIExtractor) getSecuritySchema(name string, document *openapi3.T) *openapi3.SecurityScheme {
	for name, ref := range document.Components.SecuritySchemes {
		if name == name {
			return ref.Value
		}
	}

	return nil
}
