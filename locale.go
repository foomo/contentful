package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// LocalesService service
type LocalesService service

// Locale model
type Locale struct {
	Sys *Sys `json:"sys,omitempty"`

	// Locale name
	Name string `json:"name,omitempty"`

	// Language code
	Code string `json:"code,omitempty"`

	// If no content is provided for the locale, the Delivery API will return content in a locale specified below:
	FallbackCode string `json:"fallbackCode,omitempty"`

	// Make the locale as default locale for your account
	Default bool `json:"default,omitempty"`

	// Entries with required fields can still be published if locale is empty.
	Optional bool `json:"optional,omitempty"`

	// Includes locale in the Delivery API response.
	CDA bool `json:"contentDeliveryApi"`

	// Displays locale to editors and enables it in Management API.
	CMA bool `json:"contentManagementApi"`
}

// GetVersion returns entity version
func (locale *Locale) GetVersion() int {
	version := 1
	if locale.Sys != nil {
		version = locale.Sys.Version
	}

	return version
}

// List returns a locales collection
func (service *LocalesService) List(ctx context.Context, spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s%s/locales", spaceID, getEnvPath(service.c))
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single locale entity
func (service *LocalesService) Get(ctx context.Context, spaceID, localeID string) (*Locale, error) {
	path := fmt.Sprintf("/spaces/%s%s/locales/%s", spaceID, getEnvPath(service.c), localeID)
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var locale Locale
	if err := service.c.do(req, &locale); err != nil {
		return nil, err
	}

	return &locale, nil
}

// Delete the locale
func (service *LocalesService) Delete(ctx context.Context, spaceID string, locale *Locale) error {
	path := fmt.Sprintf("/spaces/%s%s/locales/%s", spaceID, getEnvPath(service.c), locale.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(locale.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Upsert updates or creates a new locale entity
func (service *LocalesService) Upsert(ctx context.Context, spaceID string, locale *Locale) error {
	bytesArray, err := json.Marshal(locale)
	if err != nil {
		return err
	}

	var path string
	var method string

	if locale.Sys != nil && locale.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s%s/locales/%s", spaceID, getEnvPath(service.c), locale.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s%s/locales", spaceID, getEnvPath(service.c))
		method = "POST"
	}

	req, err := service.c.newRequest(ctx, method, path, nil, bytes.NewReader(bytesArray), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(locale.GetVersion()))

	return service.c.do(req, locale)
}
