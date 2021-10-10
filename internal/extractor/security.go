package extractor

import (
	"openapi-code-sample-generator/internal/log"
	"openapi-code-sample-generator/internal/types"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetSecurity returns security related parameters
func (o *openAPIExtractor) GetSecurity(operation *openapi3.Operation, document *openapi3.T) (parameters types.Parameters, basicAuth bool, err error) {
	parameters.Query = make([]*types.Parameter, 0)
	parameters.Header = make([]*types.Parameter, 0)
	parameters.Path = make([]*types.Parameter, 0)
	parameters.Cookie = make([]*types.Parameter, 0)
	basicAuth = false

	requirements := operation.Security

	if requirements != nil {
		for _, requirement := range *requirements {
			for name := range requirement {
				security := o.getSecuritySchema(name, document)
				if security == nil {
					log.Error("Not found")
					continue
				}
				switch strings.ToLower(security.Type) {
				case "http":
					switch strings.ToLower(security.Scheme) {
					case "basic":
						basicAuth = true
					case "bearer":
						parameters.Header = append(parameters.Header, &types.Parameter{
							Name:  "Authorization",
							Value: "Bearer ${TOKEN}",
						})
					}
				case "apikey":
					switch strings.ToLower(security.In) {
					case "query":
						parameters.Query = append(parameters.Query, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					case "cookie":
						parameters.Cookie = append(parameters.Cookie, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					case "header":
						parameters.Header = append(parameters.Header, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					}
				case "openidconnect", "oauth2":
					parameters.Header = append(parameters.Header, &types.Parameter{
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
	for secName, ref := range document.Components.SecuritySchemes {
		if name == secName {
			return ref.Value
		}
	}

	return nil
}
