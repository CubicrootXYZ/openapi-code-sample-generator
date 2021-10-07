package helper

import (
	"openapi-code-sample-generator/internal/types"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetSecurity returns security related parameters
func GetSecurity(operation *openapi3.Operation, document *openapi3.T) (queryParameters []*types.Parameter, headerParameters []*types.Parameter, cookieParameters []*types.Parameter, basicAuth bool) {
	queryParameters = make([]*types.Parameter, 0)
	headerParameters = make([]*types.Parameter, 0)
	basicAuth = false

	requirements := operation.Security

	if requirements != nil {
		for _, requirement := range *requirements {
			for name := range requirement {
				security := getSecuritySchema(name, document)
				if security == nil {
					continue
				}
				switch strings.ToLower(security.Type) {
				case "http":
					switch strings.ToLower(security.Scheme) {
					case "basic":
						basicAuth = true
					case "bearer":
						headerParameters = append(headerParameters, &types.Parameter{
							Name:  "Authorization",
							Value: "Bearer ${TOKEN}",
						})
					}
				case "apikey":
					switch strings.ToLower(security.In) {
					case "query":
						queryParameters = append(queryParameters, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					case "cookie":
						cookieParameters = append(cookieParameters, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					case "header":
						headerParameters = append(headerParameters, &types.Parameter{
							Name:  security.Name,
							Value: "${TOKEN}",
						})
					}
				case "openidconnect", "oauth2":
					headerParameters = append(headerParameters, &types.Parameter{
						Name:  "Authorization",
						Value: "Bearer ${TOKEN}",
					})
				}
			}
		}
	}

	return
}

func getSecuritySchema(name string, document *openapi3.T) *openapi3.SecurityScheme {
	for name, ref := range document.Components.SecuritySchemes {
		if name == name {
			return ref.Value
		}
	}

	return nil
}
