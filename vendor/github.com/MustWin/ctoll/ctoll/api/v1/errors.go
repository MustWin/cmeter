package v1

import (
	"net/http"

	"github.com/MustWin/ctoll/ctoll/api/errcode"
)

const ErrorGroup = "ctoll.api.v1"

var (
	ErrorCodeOrgUnknown = errcode.Register(ErrorGroup, errcode.ErrorDescriptor{
		Value:          "ORG_UNKNOWN",
		Message:        "organization not known to server",
		Description:    "This is returned if the organization ID used during an operation is unknown to the server.",
		HTTPStatusCode: http.StatusNotFound,
	})

	ErrorCodeAPIKeyInvalid = errcode.Register(ErrorGroup, errcode.ErrorDescriptor{
		Value:          "API_KEY_INVALID",
		Message:        "The api key was missing or is invalid.",
		Description:    "An API key is required to access the resource and one was not provided, doesn't exist, or was otherwise invalid. Verify your key and try your request again.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	ErrorCodeMeterEventInvalid = errcode.Register(ErrorGroup, errcode.ErrorDescriptor{
		Value:          "METER_EVENT_INVALID",
		Message:        "The meter event was missing or invalid.",
		Description:    "The meter event body was not provided, doesn't exist, or was otherwise invalid. Verify your entity and try your request again.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	ErrorCodeUnsupportedMeterEventType = errcode.Register(ErrorGroup, errcode.ErrorDescriptor{
		Value:          "UNSUPPORTED_METER_EVENT_TYPE",
		Message:        "The meter event type is not supported.",
		Description:    "The meter event type is unsupported by the server. Verify your entity and try your request again.",
		HTTPStatusCode: http.StatusBadRequest,
	})

	ErrorCodeBillingModelInvalid = errcode.Register(ErrorGroup, errcode.ErrorDescriptor{
		Value:          "BILLING_MODEL_INVALID",
		Message:        "The billing model is missing or invalid.",
		Description:    "The billing model body was not provided or was invalid. Verify your entity and try your request again.",
		HTTPStatusCode: http.StatusBadRequest,
	})
)
