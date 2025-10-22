package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// APIKeyService service
type APIKeyService service

// APIKey model
type APIKey struct {
	Sys           *Sys            `json:"sys,omitempty"`
	Name          string          `json:"name,omitempty"`
	Description   string          `json:"description,omitempty"`
	AccessToken   string          `json:"accessToken,omitempty"`
	Policies      []*APIKeyPolicy `json:"policies,omitempty"`
	PreviewAPIKey *PreviewAPIKey  `json:"preview_api_key,omitempty"`
}

// APIKeyPolicy model
type APIKeyPolicy struct {
	Effect  string `json:"effect,omitempty"`
	Actions string `json:"actions,omitempty"`
}

// PreviewAPIKey model
type PreviewAPIKey struct {
	Sys *Sys `json:"sys,omitempty"`
}

// MarshalJSON for custom json marshaling
func (apiKey *APIKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}{
		Name:        apiKey.Name,
		Description: apiKey.Description,
	})
}

// GetVersion returns entity version
func (apiKey *APIKey) GetVersion() int {
	version := 1
	if apiKey.Sys != nil {
		version = apiKey.Sys.Version
	}

	return version
}

// List returns all api keys collection
func (service *APIKeyService) List(ctx context.Context, spaceID string) *Collection[APIKey] {
	path := fmt.Sprintf("/spaces/%s%s/api_keys", spaceID, getEnvPath(service.c))
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return &Collection[APIKey]{}
	}

	col := NewCollection[APIKey](&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single api key entity
func (service *APIKeyService) Get(ctx context.Context, spaceID, apiKeyID string) (*APIKey, error) {
	path := fmt.Sprintf("/spaces/%s%s/api_keys/%s", spaceID, getEnvPath(service.c), apiKeyID)
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var apiKey APIKey
	if err := service.c.do(req, &apiKey); err != nil {
		return nil, err
	}

	return &apiKey, nil
}

// Upsert updates or creates a new api key entity
func (service *APIKeyService) Upsert(ctx context.Context, spaceID string, apiKey *APIKey) error {
	bytesArray, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}

	var path string
	var method string

	if apiKey.Sys != nil && apiKey.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s%s/api_keys/%s", spaceID, getEnvPath(service.c), apiKey.Sys.ID)
		method = http.MethodPut
	} else {
		path = fmt.Sprintf("/spaces/%s%s/api_keys", spaceID, getEnvPath(service.c))
		method = http.MethodPost
	}

	req, err := service.c.newRequest(ctx, method, path, nil, bytes.NewReader(bytesArray), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(apiKey.GetVersion()))

	return service.c.do(req, apiKey)
}

// Delete deletes a sinlge api key entity
func (service *APIKeyService) Delete(ctx context.Context, spaceID string, apiKey *APIKey) error {
	path := fmt.Sprintf("/spaces/%s%s/api_keys/%s", spaceID, getEnvPath(service.c), apiKey.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(apiKey.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}
