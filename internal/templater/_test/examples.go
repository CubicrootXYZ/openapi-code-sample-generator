package test

import (
	"github.com/CubicrootXYZ/openapi-code-sample-generator/internal/templater"
	"github.com/getkin/kin-openapi/openapi3"
)

func Endpoint1() *templater.Endpoint {
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
						Type: &openapi3.Types{"integer"},
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
						Type: &openapi3.Types{"integer"},
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
						Type:   &openapi3.Types{"string"},
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
						Type: &openapi3.Types{"boolean"},
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
						Type: &openapi3.Types{"array"},
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
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
						Type: &openapi3.Types{"object"},
						AnyOf: openapi3.SchemaRefs{
							nil,
							&openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"object"},
									Properties: openapi3.Schemas{"param6-sub": &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: &openapi3.Types{"array"},
											Items: &openapi3.SchemaRef{
												Value: &openapi3.Schema{
													Type: &openapi3.Types{"string"},
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

func Endpoint2() *templater.Endpoint {
	endpoint := templater.Endpoint{}
	endpoint.HTTPVerb = "POST"
	endpoint.Path = "/random_#+!§$%&/()=/path"
	endpoint.OpenAPI.Operation = &openapi3.Operation{
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: openapi3.Schemas{"param6-sub": &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: &openapi3.Types{"array"},
										Items: &openapi3.SchemaRef{
											Value: &openapi3.Schema{
												Type: &openapi3.Types{"string"},
											},
										},
									},
								}},
							},
						},
					},
				},
			},
		},
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

func Endpoint3() *templater.Endpoint {
	endpoint := templater.Endpoint{}
	endpoint.HTTPVerb = "POST"
	endpoint.Path = "/random_#+!§$%&/()=/path"
	endpoint.OpenAPI.Operation = &openapi3.Operation{
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/xml": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: openapi3.Schemas{"param6-sub": &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: &openapi3.Types{"array"},
										Items: &openapi3.SchemaRef{
											Value: &openapi3.Schema{
												Type: &openapi3.Types{"string"},
											},
										},
									},
								}},
							},
						},
					},
				},
			},
		},
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

// Endpoint4 is a simple endpoint with 2 path variables.
func Endpoint4() *templater.Endpoint {
	endpoint := templater.Endpoint{}
	endpoint.HTTPVerb = "POST"
	endpoint.Path = "/data/{id1}/{id2}"
	endpoint.OpenAPI.Operation = &openapi3.Operation{
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
							},
						},
					},
				},
			},
		},
	}
	endpoint.OpenAPI.PathItem = &openapi3.PathItem{
		Post: endpoint.OpenAPI.Operation,
		Parameters: openapi3.Parameters{
			{
				Value: &openapi3.Parameter{
					Name: "id1",
					In:   "path",
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: &openapi3.Types{
								"string",
							},
						},
					},
					Required: true,
				},
			},
			{
				Value: &openapi3.Parameter{
					Name: "id2",
					In:   "path",
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: &openapi3.Types{
								"integer",
							},
						},
					},
					Required: true,
				},
			},
		},
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
