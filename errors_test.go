package contentful

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotFoundErrorResponse(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintln(w, readTestData(t, "error-notfound.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test space
	_, err = cma.Spaces.Get(context.TODO(), "unknown-space-id")
	require.Error(t, err)
	var notFoundError NotFoundError
	ok := errors.As(err, &notFoundError)
	assert.True(t, ok)
	assert.Equal(t, 404, notFoundError.APIError.res.StatusCode)
	assert.Equal(t, "request-id", notFoundError.APIError.err.RequestID)
	assert.Equal(t, "The resource could not be found.", notFoundError.APIError.err.Message)
	assert.Equal(t, "Error", notFoundError.APIError.err.Sys.Type)
	assert.Equal(t, "NotFound", notFoundError.APIError.err.Sys.ID)
}

func TestRateLimitExceededResponse(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintln(w, readTestData(t, "error-ratelimit.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test space
	space := &Space{Name: "test-space"}
	err = cma.Spaces.Upsert(context.TODO(), space)
	require.Error(t, err)
	var rateLimitExceededError RateLimitExceededError
	ok := errors.As(err, &rateLimitExceededError)
	assert.True(t, ok)
	assert.Equal(t, 403, rateLimitExceededError.APIError.res.StatusCode)
	assert.Equal(t, "request-id", rateLimitExceededError.APIError.err.RequestID)
	assert.Equal(t, "You are creating too many Spaces.", rateLimitExceededError.APIError.err.Message)
	assert.Equal(t, "Error", rateLimitExceededError.APIError.err.Sys.Type)
	assert.Equal(t, "RateLimitExceeded", rateLimitExceededError.APIError.err.Sys.ID)
}
