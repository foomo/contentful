package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// SpacesService model
type SpacesService service

// Space model
type Space struct {
	Sys           *Sys   `json:"sys,omitempty"`
	Name          string `json:"name,omitempty"`
	DefaultLocale string `json:"defaultLocale,omitempty"`
}

// MarshalJSON for custom json marshaling
func (space *Space) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name          string `json:"name,omitempty"`
		DefaultLocale string `json:"defaultLocale,omitempty"`
	}{
		Name:          space.Name,
		DefaultLocale: space.DefaultLocale,
	})
}

// GetVersion returns entity version
func (space *Space) GetVersion() int {
	version := 1
	if space.Sys != nil {
		version = space.Sys.Version
	}

	return version
}

// List creates a spaces collection
func (service *SpacesService) List(ctx context.Context) *Collection[Space] {
	req, _ := service.c.newRequest(ctx, http.MethodGet, "/spaces", nil, nil, nil)

	col := NewCollection[Space](&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single space entity
func (service *SpacesService) Get(ctx context.Context, spaceID string) (*Space, error) {
	path := fmt.Sprintf("/spaces/%s", spaceID)
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var space *Space
	if err := service.c.do(req, &space); err != nil {
		return nil, err
	}

	return space, nil
}

// Upsert updates or creates a new space
func (service *SpacesService) Upsert(ctx context.Context, space *Space) error {
	bytesArray, err := Marshal(space)
	if err != nil {
		return err
	}

	var path string
	var method string

	if space.Sys != nil && space.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s%s", space.Sys.ID, getEnvPath(service.c))
		method = "PUT"
	} else {
		path = "/spaces"
		method = "POST"
	}

	req, err := service.c.newRequest(ctx, method, path, nil, bytes.NewReader(bytesArray), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(space.GetVersion()))

	return service.c.do(req, &space)
}

// Delete the given space
func (service *SpacesService) Delete(ctx context.Context, space *Space) error {
	path := fmt.Sprintf("/spaces/%s", space.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(space.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}
