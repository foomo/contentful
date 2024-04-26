package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	server         *httptest.Server
	cma            *Contentful
	c              *Contentful
	CMAToken       = "b4c0n73n7fu1"
	CDAToken       = "cda-token"
	CPAToken       = "cpa-token"
	spaceID        = "id1"
	organizationID = "org-id"
)

func readTestData(t *testing.T, fileName string) string {
	t.Helper()
	path := "testdata/" + fileName
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(content)
}

func checkHeaders(t *testing.T, req *http.Request) {
	t.Helper()
	assert.Equal(t, "Bearer "+CMAToken, req.Header.Get("Authorization"))
	assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
}

func spaceFromTestData(t *testing.T, fileName string) (*Space, error) {
	t.Helper()
	content := readTestData(t, fileName)

	var space Space
	err := json.NewDecoder(strings.NewReader(content)).Decode(&space)
	require.NoError(t, err)

	return &space, nil
}

func webhookFromTestData(t *testing.T, fileName string) (*Webhook, error) {
	t.Helper()
	content := readTestData(t, fileName)

	var webhook Webhook
	err := json.NewDecoder(strings.NewReader(content)).Decode(&webhook)
	require.NoError(t, err)

	return &webhook, nil
}

func contentTypeFromTestData(t *testing.T, fileName string) (*ContentType, error) {
	t.Helper()
	content := readTestData(t, fileName)

	var ct ContentType
	err := json.NewDecoder(strings.NewReader(content)).Decode(&ct)
	require.NoError(t, err)

	return &ct, nil
}

func localeFromTestData(t *testing.T, fileName string) (*Locale, error) {
	t.Helper()
	content := readTestData(t, fileName)

	var locale Locale
	err := json.NewDecoder(strings.NewReader(content)).Decode(&locale)
	require.NoError(t, err)

	return &locale, nil
}

func setup() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture := strings.ReplaceAll(r.URL.Path, "/", "-")
		fixture = strings.TrimLeft(fixture, "-")
		var path string

		if e := r.URL.Query().Get("error"); e != "" {
			path = "testdata/error-" + e + ".json"
		} else {
			if r.Method == http.MethodGet {
				path = "testdata/" + fixture + ".json"
			}

			if r.Method == http.MethodPost {
				path = "testdata/" + fixture + "-new.json"
			}

			if r.Method == http.MethodPut {
				path = "testdata/" + fixture + "-updated.json"
			}
		}

		file, err := os.ReadFile(path)
		if err != nil {
			_, _ = fmt.Fprintln(w, err)
			return
		}

		_, _ = fmt.Fprintln(w, string(file))
	})

	server = httptest.NewServer(handler)

	c = NewCMA(CMAToken)
	c.BaseURL = server.URL
}

func teardown() {
	server.Close()
	c = nil
}

func TestContentfulNewCMA(t *testing.T) {
	cma := NewCMA(CMAToken)
	assert.IsType(t, Contentful{}, *cma)
	assert.Equal(t, "https://api.contentful.com", cma.BaseURL)
	assert.Equal(t, "CMA", cma.api)
	assert.Equal(t, CMAToken, cma.token)
	assert.Equal(t, fmt.Sprintf("Bearer %s", CMAToken), cma.Headers["Authorization"])
	assert.Equal(t, "application/vnd.contentful.management.v1+json", cma.Headers["Content-Type"])
	assert.Equal(t, fmt.Sprintf("sdk contentful.go/%s", Version), cma.Headers["X-Contentful-User-Agent"])
}

func TestContentfulNewCDA(t *testing.T) {
	cda := NewCDA(CDAToken)
	assert.IsType(t, Contentful{}, *cda)
	assert.Equal(t, "https://cdn.contentful.com", cda.BaseURL)
	assert.Equal(t, "CDA", cda.api)
	assert.Equal(t, CDAToken, cda.token)
	assert.Equal(t, fmt.Sprintf("Bearer %s", CDAToken), cda.Headers["Authorization"])
	assert.Equal(t, "application/vnd.contentful.delivery.v1+json", cda.Headers["Content-Type"])
	assert.Equal(t, fmt.Sprintf("contentful-go/%s", Version), cda.Headers["X-Contentful-User-Agent"])
}

func TestContentfulNewCPA(t *testing.T) {
	cpa := NewCPA(CPAToken)
	assert.IsType(t, Contentful{}, *cpa)
	assert.Equal(t, "https://preview.contentful.com", cpa.BaseURL)
	assert.Equal(t, "CPA", cpa.api)
	assert.Equal(t, CPAToken, cpa.token)
}

func TestContentfulSetOrganization(t *testing.T) {
	cma := NewCMA(CMAToken)
	cma.SetOrganization(organizationID)
	assert.Equal(t, organizationID, cma.Headers["X-Contentful-Organization"])
}

func TestContentfulSetClient(t *testing.T) {
	newClient := &http.Client{}
	cma := NewCMA(CMAToken)
	cma.SetHTTPClient(newClient)
	assert.Exactly(t, newClient, cma.client)
}

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	method := http.MethodGet
	path := "/some/path"
	query := url.Values{}
	query.Add("foo", "bar")
	query.Add("faz", "zoo")

	expectedURL, _ := url.Parse(c.BaseURL)
	expectedURL.Path = path
	expectedURL.RawQuery = query.Encode()

	req, err := c.newRequest(context.TODO(), method, path, query, nil)
	require.NoError(t, err)
	assert.Equal(t, "Bearer "+CMAToken, req.Header.Get("Authorization"))
	assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
	assert.Equal(t, req.Method, method)
	assert.Equal(t, req.URL.String(), expectedURL.String())

	method = "POST"
	type RequestBody struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	bodyData := RequestBody{
		Name: "test",
		Age:  10,
	}
	body, err := json.Marshal(bodyData)
	require.NoError(t, err)
	req, err = c.newRequest(context.TODO(), method, path, query, bytes.NewReader(body))
	require.NoError(t, err)
	assert.Equal(t, "Bearer "+CMAToken, req.Header.Get("Authorization"))
	assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
	assert.Equal(t, req.Method, method)
	assert.Equal(t, req.URL.String(), expectedURL.String())
	defer req.Body.Close()
	var requestBody RequestBody
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	require.NoError(t, err)
	assert.Equal(t, requestBody, bodyData)
}

func TestHandleError(t *testing.T) {
	setup()
	defer teardown()

	method := http.MethodGet
	path := "/some/path"
	requestID := "request-id"
	query := url.Values{}
	errResponse := ErrorResponse{
		Sys: &Sys{
			ID:   "AccessTokenInvalid",
			Type: "Error",
		},
		Message:   "Access token is invalid",
		RequestID: requestID,
	}

	marshaled, err := json.Marshal(errResponse)
	require.NoError(t, err)
	errResponseReader := bytes.NewReader(marshaled)
	errResponseReadCloser := io.NopCloser(errResponseReader)

	req, err := c.newRequest(context.TODO(), method, path, query, nil)
	require.NoError(t, err)
	responseHeaders := http.Header{}
	responseHeaders.Add("X-Contentful-Request-Id", requestID)
	res := &http.Response{
		Header:     responseHeaders,
		StatusCode: http.StatusUnauthorized,
		Body:       errResponseReadCloser,
		Request:    req,
	}

	err = c.handleError(req, res)
	assert.IsType(t, AccessTokenInvalidError{}, err)
	assert.Equal(t, req, err.(AccessTokenInvalidError).APIError.req)          //nolint:errorlint
	assert.Equal(t, res, err.(AccessTokenInvalidError).APIError.res)          //nolint:errorlint
	assert.Equal(t, &errResponse, err.(AccessTokenInvalidError).APIError.err) //nolint:errorlint
}

func TestBackoffForPerSecondLimiting(t *testing.T) {
	var err error
	var backoff atomic.Bool
	waitSeconds := 2

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if backoff.Load() {
			w.Header().Set("X-Contentful-Request-Id", "request-id")
			w.Header().Set("Content-Type", "application/vnd.contentful.management.v1+json")
			w.Header().Set("X-Contentful-Ratelimit-Hour-Limit", "36000")
			w.Header().Set("X-Contentful-Ratelimit-Hour-Remaining", "35883")
			w.Header().Set("X-Contentful-Ratelimit-Reset", strconv.Itoa(waitSeconds))
			w.Header().Set("X-Contentful-Ratelimit-Second-Limit", "10")
			w.Header().Set("X-Contentful-Ratelimit-Second-Remaining", "0")
			w.WriteHeader(http.StatusTooManyRequests)

			_, _ = w.Write([]byte(readTestData(t, "error-ratelimit.json")))
		} else {
			_, _ = w.Write([]byte(readTestData(t, "space-1.json")))
		}
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	go func() {
		time.Sleep(time.Second * time.Duration(waitSeconds))
		backoff.Store(true)
	}()

	space, err := cma.Spaces.Get(context.TODO(), "id1")
	require.NoError(t, err)
	assert.Equal(t, "Contentful Example API", space.Name)
	assert.Equal(t, "id1", space.Sys.ID)
}
