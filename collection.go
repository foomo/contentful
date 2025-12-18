package contentful

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
)

// CollectionOptions holds init options
type CollectionOptions struct {
	Limit uint16
}

// Includes model
type Includes struct {
	Entry map[string]interface{} `json:"Entry"`
	Asset map[string]interface{} `json:"Asset"`
}

// Collection model with generic type parameter for Items
type Collection[T any] struct {
	Query
	c           *Contentful
	req         *http.Request
	page        uint16
	Sys         *Sys                   `json:"sys"`
	Total       int                    `json:"total"`
	Skip        int                    `json:"skip"`
	Limit       uint16                 `json:"limit"`
	Items       []T                    `json:"items"`
	Includes    map[string]interface{} `json:"includes"`
	NextSyncURL string                 `json:"nextSyncUrl"`
	NextPageURL string                 `json:"nextPageUrl"`
	SyncToken   string                 `json:"syncToken"`
	// Errors which occur in the contentful structure. They are not checked in
	// this source code. Please do it yourself as you might still want to parse
	// the result despite the error.
	Errors []Error `json:"errors"`
	// Details might also get set in case of errors.
	Details *ErrorDetails `json:"details"`
}

var syncTokenRegex = regexp.MustCompile(`sync_token=([a-zA-Z0-9_\-]+)`)

// NewCollection initializes a new generic collection
func NewCollection[T any](options *CollectionOptions) *Collection[T] {
	query := NewQuery()
	limit := uint16(100)
	if options.Limit > 0 {
		limit = options.Limit
	}
	query.Limit(limit)

	return &Collection[T]{
		Query: *query,
		page:  1,
		Limit: limit,
	}
}

// Next makes the col.req
func (col *Collection[T]) Next() (*Collection[T], error) {
	// setup query params
	if col.SyncToken != "" {
		col.Query = *NewQuery()
		col.Query.SyncToken(col.SyncToken)
	} else {
		skip := col.limit * (col.page - 1)
		col.Query.Skip(skip)
	}

	// override request query
	col.req.URL.RawQuery = col.String()
	col.Sys = nil
	col.Items = nil
	col.Includes = nil
	col.Errors = nil
	col.Details = nil

	// makes api call
	err := col.c.do(col.req, &col)
	if err != nil {
		return nil, err
	}

	col.page++
	if col.NextPageURL != "" {
		syncToken := syncTokenRegex.FindStringSubmatch(col.NextPageURL)
		col.SyncToken = syncToken[1]
	} else if col.NextSyncURL != "" {
		syncToken := syncTokenRegex.FindStringSubmatch(col.NextSyncURL)
		col.SyncToken = syncToken[1]
	}
	return col, nil
}

// Get makes the col.req with no automatic pagination
func (col *Collection[T]) Get() (*Collection[T], error) {
	// override request query
	col.req.URL.RawQuery = col.String()
	// makes api call
	err := col.c.do(col.req, &col)
	if err != nil {
		return nil, err
	}

	return col, nil
}

// GetAll paginates and returns all items - beware of memory usage!
func (col *Collection[T]) GetAll() (*Collection[T], error) {
	var allItems []T
	col.Query.Limit(col.Limit)
	for {
		var errNext error
		col, errNext = col.Next()
		if errNext != nil {
			return nil, errNext
		}
		allItems = append(allItems, col.Items...)
		if len(col.Items) < int(col.Limit) {
			break
		}
	}
	col.Items = allItems
	return col, nil
}

// ToIncludesEntry cast includesEntry to Entry model
func (col *Collection[T]) ToIncludesEntry() ([]*Entry, error) {
	var includesEntry []*Entry

	byteArray, err := json.Marshal(col.Includes["Entry"])
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesEntry); err != nil {
		return nil, err
	}

	return includesEntry, nil
}

// ToIncludesEntryMap returns a map of Entry's from the Includes
func (col *Collection[T]) ToIncludesEntryMap() (map[string]*Entry, error) {
	var includesEntry []*Entry
	includesEntryMap := make(map[string]*Entry)

	byteArray, err := json.Marshal(col.Includes["Entry"])
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesEntry); err != nil {
		return nil, err
	}

	for _, e := range includesEntry {
		includesEntryMap[e.Sys.ID] = e
	}
	return includesEntryMap, nil
}

// ToIncludesAsset cast includesAsset to Asset model
func (col *Collection[T]) ToIncludesAsset() ([]*IncludeAsset, error) {
	var includesAsset []*IncludeAsset

	byteArray, err := json.Marshal(col.Includes["Asset"])
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesAsset); err != nil {
		return nil, err
	}
	return includesAsset, nil
}

// ToIncludesAssetMap returns a map of Asset's from the Includes
func (col *Collection[T]) ToIncludesAssetMap() (map[string]*IncludeAsset, error) {
	var includesAsset []*IncludeAsset
	includesAssetMap := make(map[string]*IncludeAsset)

	byteArray, err := json.Marshal(col.Includes["Asset"])
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesAsset); err != nil {
		return nil, err
	}

	for _, a := range includesAsset {
		includesAssetMap[a.Sys.ID] = a
	}
	return includesAssetMap, nil
}

// ToIncludesLocalizedAssetMap returns a map of Asset's from the Includes
func (col *Collection[T]) ToIncludesLocalizedAssetMap() (map[string]*Asset, error) {
	var includesAsset []*Asset
	includesAssetMap := make(map[string]*Asset)

	byteArray, err := json.Marshal(col.Includes["Asset"])
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesAsset); err != nil {
		return nil, err
	}

	for _, a := range includesAsset {
		includesAssetMap[a.Sys.ID] = a
	}
	return includesAssetMap, nil
}
