package contentful

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// TagsService servıce
type TagsService service

// Tag model
type Tag struct {
	Sys  *Sys   `json:"sys"`
	Name string `json:"name,omitempty"`
}

// List returns tags collection
func (service *TagsService) List(ctx context.Context, spaceID string) *Collection[Tag] {
	path := fmt.Sprintf("/spaces/%s%s/tags", spaceID, getEnvPath(service.c))
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return &Collection[Tag]{}
	}

	col := NewCollection[Tag](&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single entry
func (service *TagsService) Get(ctx context.Context, spaceID, tagID string, locale ...string) (*Tag, error) {
	path := fmt.Sprintf("/spaces/%s%s/entries/%s", spaceID, getEnvPath(service.c), tagID)
	query := url.Values{}
	if len(locale) > 0 {
		query["locale"] = locale
	}
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, query, nil, nil)
	if err != nil {
		return &Tag{}, err
	}

	var tag *Tag
	if err := service.c.do(req, &tag); err != nil {
		return nil, err
	}

	return tag, err
}
