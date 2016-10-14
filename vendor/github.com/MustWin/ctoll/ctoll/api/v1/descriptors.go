package v1

import (
	"net/http"
	"regexp"

	"github.com/MustWin/ctoll/ctoll/api/describe"
	"github.com/MustWin/ctoll/ctoll/api/errcode"
)

var (
	IDRegex = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[1][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`)

	apiKeyParameter = describe.ParameterDescriptor{
		Name:        "api_key",
		Type:        "string",
		Format:      IDRegex.String(),
		Required:    true,
		Description: "An API key.",
	}

	apiKeyHeader = describe.ParameterDescriptor{
		Name:        "cToll-Api-Key",
		Type:        "string",
		Description: "The API key to register the request under.",
		Format:      IDRegex.String(),
		Examples:    []string{"2390511a-870d-11e6-ae22-56b6b6499612"},
	}

	versionHeader = describe.ParameterDescriptor{
		Name:        "cToll-Api-Version",
		Type:        "string",
		Description: "The build version of the cToll API server.",
		Format:      "<version>",
		Examples:    []string{"0.0.0-dev"},
	}

	hostHeader = describe.ParameterDescriptor{
		Name:        "Host",
		Type:        "string",
		Description: "",
		Format:      "<hostname>",
		Examples:    []string{"api.ctoll.io"},
	}

	orgIDParameter = describe.ParameterDescriptor{
		Name:        "org_id",
		Type:        "string",
		Description: "Identifier for organization",
		Format:      IDRegex.String(),
		Required:    true,
	}

	jsonContentLengthHeader = describe.ParameterDescriptor{
		Name:        "Content-Length",
		Type:        "integer",
		Description: "Length of the JSON body.",
		Format:      "<length>",
	}

	zeroContentLengthHeader = describe.ParameterDescriptor{
		Name:        "Content-Length",
		Type:        "integer",
		Description: "The 'Content-Length' header must be zero and the body must be empty.",
		Format:      "0",
	}

	apiKeyInvalidResp = describe.ResponseDescriptor{
		Name:        "Invalid API Key",
		Description: "The API key was missing or invalid.",
		StatusCode:  http.StatusBadRequest,
		Headers: []describe.ParameterDescriptor{
			versionHeader,
		},
		Body: describe.BodyDescriptor{
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		},
		ErrorCodes: []errcode.ErrorCode{
			ErrorCodeAPIKeyInvalid,
		},
	}

	orgNotFoundResp = describe.ResponseDescriptor{
		Name:        "Organization Unknown Error",
		StatusCode:  http.StatusNotFound,
		Description: "The organization is not known to the server.",
		Headers: []describe.ParameterDescriptor{
			versionHeader,
			jsonContentLengthHeader,
		},
		Body: describe.BodyDescriptor{
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		},
		ErrorCodes: []errcode.ErrorCode{
			ErrorCodeOrgUnknown,
		},
	}
)

var (
	errorsBody = `{
	"errors:" [
	    {
            "code": <error code>,
            "message": <error message>,
            "detail": ...
        },
        ...
    ]
}`

	orgBody = `{
	"id": <uuid>,
	"name": ...
}`

	orgsBody = `[
` + orgBody + `, ...
]`

	apiKeyBody = `{
	"key": <uuid>,
	"org_id": <org uuid>
}`

	apiKeyArrBody = `[
	` + apiKeyBody + `, ...
]`

	billingModelBody = `{
	"method": <pricing method>,
	"usage": {
		"vcpu": <unit cost>,
		"memory_mb": <unit cost>,
		"disk_io": <unit cost>,
		"net_io": <unit cost>
	}, 
	"block": {
		"vcpu": <unit cost>,
		"memory_mb": <unit cost>
	}
}`
)

var APIDescriptor = struct {
	RouteDescriptors []describe.RouteDescriptor `json:"routes"`
}{
	RouteDescriptors: routeDescriptors,
}

var routeDescriptors = []describe.RouteDescriptor{
	{
		Name:        RouteNameBase,
		Path:        "/v1",
		Entity:      "Base",
		Description: "Base V1 API route, can be used for lightweight health and version check.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Check that the server supports the cToll V1 API.",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "The API implements the V1 protocol and is accessible.",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},

						Failures: []describe.ResponseDescriptor{
							{
								Description: "The API does not support the V1 protocol.",
								StatusCode:  http.StatusNotFound,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameMeterEvents,
		Path:        "/v1/meter/events",
		Entity:      "",
		Description: "The main submission route for client meter events.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "POST",
				Description: "Submit a meter event to the API.",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							apiKeyHeader,
							hostHeader,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "Meter event is being processed.",
								StatusCode:  http.StatusAccepted,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},
						Failures: []describe.ResponseDescriptor{
							apiKeyInvalidResp,
							{
								Name:        "Invalid Meter Event",
								Description: "The meter event request body was invalid in some way as described by the error codes. The client should resolve the issue and retry the request.",
								StatusCode:  http.StatusBadRequest,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
								},
								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      errorsBody,
								},
								ErrorCodes: []errcode.ErrorCode{
									ErrorCodeMeterEventInvalid,
									ErrorCodeUnsupportedMeterEventType},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameOrgs,
		Path:        "/v1/orgs",
		Entity:      "[]Org",
		Description: "Base V1 API route, can be used for lightweight health and version check.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get all organizations",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "All organizations returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      orgsBody,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameMeterImageNames,
		Path:        "/v1/orgs/{org_id:" + IDRegex.String() + "}/image_names",
		Entity:      "[]string",
		Description: "Route to retrieve image names belonging to a specific organization.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get all image names for the organization",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "All image names returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      "",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameMeterImageTags,
		Path:        "/v1/orgs/{org_id:" + IDRegex.String() + "}/image_tags",
		Entity:      "[]string",
		Description: "Route to retrieve image tags belonging to a specific organization.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get all image tags for the organization",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "All image tags returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      "",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameMeterLabels,
		Path:        "/v1/orgs/{org_id:" + IDRegex.String() + "}/labels",
		Entity:      "[]string",
		Description: "Route to retrieve labels belonging to a specific organization.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get all labels for the organization",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "All labels returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      "",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameMeterSessions,
		Path:        "/v1/orgs/{org_id:" + IDRegex.String() + "}/sessions",
		Entity:      "[]MeterSession",
		Description: "Route to search and retrieve meter sessions belonging to a specific organization.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get all meter sessions for the organization",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						QueryParameters: []describe.ParameterDescriptor{
							{
								Name:        "state",
								Type:        "string",
								Description: "Session State to filter by",
								Format:      "<active|closed>",
							},
							{
								Name:        "start",
								Type:        "int64",
								Description: "Session start time",
								Format:      "<unixtimestamp>",
							},
							{
								Name:        "end",
								Type:        "int64",
								Description: "Session end time",
								Format:      "<unixtimestamp>",
							},
							{
								Name:        "image_name",
								Type:        "string",
								Description: "Container image name",
								Format:      "<name>",
							},
							{
								Name:        "image_tag",
								Type:        "string",
								Description: "Container image tag",
								Format:      "<name>",
							},
							{
								Name:        "image_label",
								Type:        "[]string",
								Description: "Container image labels",
								Format:      "<label[,label,...]>",
							},
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "All meter sessions returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      "",
								},
							},
						},

						Failures: []describe.ResponseDescriptor{
							orgNotFoundResp,
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameAPIKeys,
		Path:        "/v1/orgs/{org_id:" + IDRegex.String() + "}/apikeys",
		Entity:      "[]APIKey",
		Description: "Route retrieve API keys belonging to a specific organization.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get all API keys for the organization",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "All API keys returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      apiKeyArrBody,
								},
							},
						},

						Failures: []describe.ResponseDescriptor{
							orgNotFoundResp,
						},
					},
				},
			},
			{
				Method:      "PUT",
				Description: "Create an API key for an organization.",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "API Key created",
								StatusCode:  http.StatusCreated,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      apiKeyBody,
								},
							},
						},

						Failures: []describe.ResponseDescriptor{
							orgNotFoundResp,
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameBillingModel,
		Path:        "/v1/orgs/{org_id:" + IDRegex.String() + "}/billing",
		Entity:      "BillingModel",
		Description: "Route to retrieve and modify the billing model for an organization.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get the billing model for the organization.",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "Billing model returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.BodyDescriptor{
									ContentType: "application/json; charset=utf-8",
									Format:      billingModelBody,
								},
							},
						},

						Failures: []describe.ResponseDescriptor{
							orgNotFoundResp,
						},
					},
				},
			},
			{
				Method:      "POST",
				Description: "Modify the billing model for the organization.",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "Billing model saved",
								StatusCode:  http.StatusNoContent,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},

						Failures: []describe.ResponseDescriptor{
							orgNotFoundResp,
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameAPIKey,
		Path:        "/v1/orgs/{org_id:" + IDRegex.String() + "}/apikeys/{api_key:" + IDRegex.String() + "}",
		Entity:      "APIKey",
		Description: "Route to view and perform operations on a single, specific api key.",
		Methods: []describe.MethodDescriptor{
			{
				Method:      "DELETE",
				Description: "Remove an API key from an organization.",
				Requests: []describe.RequestDescriptor{
					{
						Headers: []describe.ParameterDescriptor{
							hostHeader,
						},

						PathParameters: []describe.ParameterDescriptor{
							orgIDParameter,
							apiKeyParameter,
						},

						Successes: []describe.ResponseDescriptor{
							{
								Description: "The API key was removed successfully.",
								StatusCode:  http.StatusNoContent,
								Headers: []describe.ParameterDescriptor{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},

						Failures: []describe.ResponseDescriptor{
							orgNotFoundResp,
						},
					},
				},
			},
		},
	},
}

var routeDescriptorsMap map[string]describe.RouteDescriptor

func init() {
	routeDescriptorsMap = make(map[string]describe.RouteDescriptor, len(routeDescriptors))
	for _, descriptor := range routeDescriptors {
		routeDescriptorsMap[descriptor.Name] = descriptor
	}
}
