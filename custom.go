package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// CustomService servıce
type CustomService[T any] service

func NewCustomService[T any](c *Contentful) *CustomService[T] {
	return &CustomService[T]{c: c}
}

// GetEntryKey returns the entry's keys
func (service *CustomService[T]) GetEntryKey(ctx context.Context, entry T, key string) (*EntryField, error) {
	var base Entry
	if err := DeepCopy(&base, entry); err != nil {
		return nil, err
	}

	ef := EntryField{
		value: base.Fields[key],
	}

	col, err := service.c.ContentTypes.List(ctx, base.Sys.Space.Sys.ID).Next()
	if err != nil {
		return nil, err
	}

	for _, ct := range col.Items {
		if ct.Sys.ID != base.Sys.ContentType.Sys.ID {
			continue
		}

		for _, field := range ct.Fields {
			if field.ID != key {
				continue
			}

			ef.dataType = field.Type
		}
	}

	return &ef, nil
}

// List returns entries collection
func (service *CustomService[T]) List(ctx context.Context, spaceID string) *Collection[T] {
	path := fmt.Sprintf("/spaces/%s%s/entries", spaceID, getEnvPath(service.c))
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return &Collection[T]{}
	}

	col := NewCollection[T](&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Sync returns entries collection
func (service *CustomService[T]) Sync(ctx context.Context, spaceID string, initial bool, syncToken ...string) *Collection[T] {
	path := fmt.Sprintf("/spaces/%s%s/sync", spaceID, getEnvPath(service.c))
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return &Collection[T]{}
	}

	col := NewCollection[T](&CollectionOptions{})
	if initial {
		col.Initial("true")
	}
	if len(syncToken) == 1 {
		col.SyncToken = syncToken[0]
	}
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single entry
func (service *CustomService[T]) Get(ctx context.Context, spaceID, entryID string, locale ...string) (T, error) {
	var entry T
	path := fmt.Sprintf("/spaces/%s%s/entries/%s", spaceID, getEnvPath(service.c), entryID)
	query := url.Values{}
	if len(locale) > 0 {
		query["locale"] = locale
	}
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, query, nil, nil)
	if err != nil {
		return entry, err
	}

	if err := service.c.do(req, &entry); err != nil {
		return entry, err
	}

	return entry, err
}

// Delete the entry
func (service *CustomService[T]) Delete(ctx context.Context, spaceID string, entryID string) error {
	path := fmt.Sprintf("/spaces/%s%s/entries/%s", spaceID, getEnvPath(service.c), entryID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Upsert updates or creates a new entry
func (service *CustomService[T]) Upsert(ctx context.Context, spaceID string, entry T) error {
	var base Entry
	if err := DeepCopy(&base, entry); err != nil {
		return err
	}

	fieldsOnly := map[string]interface{}{
		"fields": base.Fields,
	}

	bytesArray, err := json.Marshal(fieldsOnly)
	if err != nil {
		return err
	}

	// Creating/updating an entry requires a content type to be provided
	if base.Sys.ContentType == nil {
		return fmt.Errorf("creating/updating an entry requires a content type")
	}

	var path string
	var method string

	if base.Sys != nil && base.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s%s/entries/%s", spaceID, getEnvPath(service.c), base.Sys.ID)
		method = http.MethodPut
	} else {
		path = fmt.Sprintf("/spaces/%s%s/entries", spaceID, getEnvPath(service.c))
		method = http.MethodPost
	}

	req, err := service.c.newRequest(ctx, method, path, nil, bytes.NewReader(bytesArray), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(base.GetVersion()))
	req.Header.Set("X-Contentful-Content-Type", base.Sys.ContentType.Sys.ID)

	return service.c.do(req, entry)
}

// Publish the entry
func (service *CustomService[T]) Publish(ctx context.Context, spaceID string, entry T) error {
	var base Entry
	if err := DeepCopy(&base, entry); err != nil {
		return err
	}
	path := fmt.Sprintf("/spaces/%s%s/entries/%s/published", spaceID, getEnvPath(service.c), base.Sys.ID)
	method := http.MethodPut

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(base.GetVersion()))

	return service.c.do(req, nil)
}

// Unpublish the entry
func (service *CustomService[T]) Unpublish(ctx context.Context, spaceID string, entry T) error {
	var base Entry
	if err := DeepCopy(&base, entry); err != nil {
		return err
	}

	path := fmt.Sprintf("/spaces/%s%s/entries/%s/published", spaceID, getEnvPath(service.c), base.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(base.GetVersion()))

	return service.c.do(req, nil)
}

// Publish the entry
func (service *CustomService[T]) Archive(ctx context.Context, spaceID string, entry T) error {
	var base Entry
	if err := DeepCopy(&base, entry); err != nil {
		return err
	}

	path := fmt.Sprintf("/spaces/%s%s/entries/%s/archived", spaceID, getEnvPath(service.c), base.Sys.ID)
	method := http.MethodPut

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(base.GetVersion()))

	return service.c.do(req, nil)
}
