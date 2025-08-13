package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse model
type ErrorResponse struct {
	Sys       *Sys          `json:"sys"`
	Message   string        `json:"message,omitempty"`
	RequestID string        `json:"requestId,omitempty"`
	Details   *ErrorDetails `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	byt, _ := json.Marshal(e)
	return string(byt)
}

// Error defines an internal contentful data model error when e.g. an assets can
// be resolved. The structure gets returned while this error is set too. The
// programmer can decide how to handle the response by checking the
// Collection.Errors field.
type Error struct {
	Sys     *Sys
	Details map[string]string
}

// ErrorDetails model
type ErrorDetails struct {
	Errors []*ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail model
type ErrorDetail struct {
	ID      string      `json:"id,omitempty"`
	Name    string      `json:"name,omitempty"`
	Path    interface{} `json:"path,omitempty"`
	Details string      `json:"details,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

// APIError model
type APIError struct {
	req *http.Request
	res *http.Response
	err *ErrorResponse
}

// AccessTokenInvalidError for 401 errors
type AccessTokenInvalidError struct {
	APIError
}

func (e AccessTokenInvalidError) Error() string {
	return e.APIError.err.Message
}

// VersionMismatchError for 409 errors
type VersionMismatchError struct {
	APIError
}

func (e VersionMismatchError) Error() string {
	return "Version " + e.APIError.req.Header.Get("X-Contentful-Version") + " is mismatched"
}

// ValidationFailedError model
type ValidationFailedError struct {
	APIError
}

func (e ValidationFailedError) Error() string {
	msg := bytes.Buffer{}

	for _, err := range e.APIError.err.Details.Errors {
		switch err.Name {
		case "uniqueFieldIds", "uniqueFieldApiNames":
			return msg.String()
		case "notResolvable":
			if err.Path != nil {
				switch path := err.Path.(type) {
				case []interface{}:
					pathString := ""
					for _, segment := range path {
						switch segment.(type) {
						case string:
							pathString += fmt.Sprintf("/%s", segment)
						}
					}
					_, _ = msg.WriteString(fmt.Sprintf("errorName: %s, path: %s", err.Name, pathString))
				}
			}
		default:
			_, _ = msg.WriteString(fmt.Sprintf("%s\n", err.Details))
		}
	}

	return msg.String()
}

// NotFoundError for 404 errors
type NotFoundError struct {
	APIError
}

func (e NotFoundError) Error() string {
	return "the requested resource can not be found"
}

// RateLimitExceededError for rate limit errors
type RateLimitExceededError struct {
	APIError
}

func (e RateLimitExceededError) Error() string {
	return e.APIError.err.Message
}

// BadRequestError error model for bad request responses
type BadRequestError struct{}

// InvalidQueryError error model for invalid query responses
type InvalidQueryError struct{}

// AccessDeniedError error model for access denied responses
type AccessDeniedError struct{}

// ServerError error model for server error responses
type ServerError struct{}
