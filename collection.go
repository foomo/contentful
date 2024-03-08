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

// Collection model
type Collection struct {
	Query
	c           *Contentful
	req         *http.Request
	page        uint16
	Sys         *Sys                   `json:"sys"`
	Total       int                    `json:"total"`
	Skip        int                    `json:"skip"`
	Limit       uint16                 `json:"limit"`
	Items       []interface{}          `json:"items"`
	Includes    map[string]interface{} `json:"includes"`
	NextSyncURL string                 `json:"nextSyncUrl"`
	NextPageURL string                 `json:"nextPageUrl"`
	SyncToken   string
	// Errors which occur in the contentful structure. They are not checked in
	// this source code. Please do it yourself as you might still want to parse
	// the result despite the error.
	Errors []Error `json:"errors"`
	// Details might also get set in case of errors.
	Details *ErrorDetails
}

// NewCollection initializes a new collection
func NewCollection(options *CollectionOptions) *Collection {
	query := NewQuery()
	limit := uint16(100)
	if options.Limit > 0 {
		limit = options.Limit
	}
	query.Limit(limit)

	return &Collection{
		Query: *query,
		page:  1,
		Limit: limit,
	}
}

// Next makes the col.req
func (col *Collection) Next() (*Collection, error) {
	// setup query params
	if col.SyncToken != "" {
		col.Query = *NewQuery()
		col.Query.SyncToken(col.SyncToken)
	} else {
		skip := col.Query.limit * (col.page - 1)
		col.Query.Skip(skip)
	}

	// override request query
	col.req.URL.RawQuery = col.Query.String()

	// makes api call
	err := col.c.do(col.req, col)
	if err != nil {
		return nil, err
	}

	col.page++
	r, _ := regexp.Compile("sync_token=([a-zA-Z0-9\\-\\_]+)")
	if col.NextPageURL != "" {
		syncToken := r.FindStringSubmatch(col.NextPageURL)
		col.SyncToken = syncToken[1]
	} else if col.NextSyncURL != "" {
		syncToken := r.FindStringSubmatch(col.NextSyncURL)
		col.SyncToken = syncToken[1]
	}
	return col, nil
}

// Get makes the col.req with no automatic pagination
func (col *Collection) Get() (*Collection, error) {
	// override request query
	col.req.URL.RawQuery = col.Query.String()
	// makes api call
	err := col.c.do(col.req, col)
	if err != nil {
		return nil, err
	}

	return col, nil
}

// GetAll paginates and returns all items - beware of memory usage!
func (col *Collection) GetAll() (*Collection, error) {
	var allItems []interface{}
	col.Query.Limit(col.Limit)
	for {
		var errNext error
		col, errNext = col.Next()
		if errNext != nil {
			return nil, errNext
		}
		allItems = append(allItems, col.Items...)
		if uint16(len(col.Items)) < col.Limit {
			break
		}
	}
	col.Items = allItems
	return col, nil
}

// ToContentType cast Items to ContentType model
func (col *Collection) ToContentType() []*ContentType {
	var contentTypes []*ContentType

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&contentTypes)

	return contentTypes
}

// ToSpace cast Items to Space model
func (col *Collection) ToSpace() []*Space {
	var spaces []*Space

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&spaces)

	return spaces
}

// ToScheduledAction cast Items to ScheduledActions model
func (col *Collection) ToScheduledAction() []*ScheduledActions {
	var scheduledActions []*ScheduledActions

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&scheduledActions)

	return scheduledActions
}

// ToEntry cast Items to Entry model
func (col *Collection) ToEntry() []*Entry {
	var entries []*Entry

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&entries)

	return entries
}

// ToLocale cast Items to Locale model
func (col *Collection) ToLocale() []*Locale {
	var locales []*Locale

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&locales)

	return locales
}

// ToAsset cast Items to Asset model
func (col *Collection) ToAsset() []*Asset {
	var assets []*Asset

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&assets)

	return assets
}

// ToAPIKey cast Items to APIKey model
func (col *Collection) ToAPIKey() []*APIKey {
	var apiKeys []*APIKey

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&apiKeys)

	return apiKeys
}

// ToWebhook cast Items to Webhook model
func (col *Collection) ToWebhook() []*Webhook {
	var webhooks []*Webhook

	byteArray, _ := json.Marshal(col.Items)
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&webhooks)

	return webhooks
}

// ToIncludesEntry cast includesEntry to Entry model
func (col *Collection) ToIncludesEntry() []*Entry {
	var includesEntry []*Entry

	byteArray, _ := json.Marshal(col.Includes["Entry"])
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesEntry)
	return includesEntry
}

// ToIncludesEntryMap returns a map of Entry's from the Includes
func (col *Collection) ToIncludesEntryMap() map[string]*Entry {
	var includesEntry []*Entry
	includesEntryMap := make(map[string]*Entry)

	byteArray, _ := json.Marshal(col.Includes["Entry"])
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesEntry)
	for _, e := range includesEntry {
		includesEntryMap[e.Sys.ID] = e
	}
	return includesEntryMap
}

// ToIncludesAsset cast includesAsset to Asset model
func (col *Collection) ToIncludesAsset() []*IncludeAsset {
	var includesAsset []*IncludeAsset

	byteArray, _ := json.Marshal(col.Includes["Asset"])
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesAsset)
	return includesAsset
}

// ToIncludesAssetMap returns a map of Asset's from the Includes
func (col *Collection) ToIncludesAssetMap() map[string]*IncludeAsset {
	var includesAsset []*IncludeAsset
	includesAssetMap := make(map[string]*IncludeAsset)

	byteArray, _ := json.Marshal(col.Includes["Asset"])
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesAsset)

	for _, a := range includesAsset {
		includesAssetMap[a.Sys.ID] = a
	}
	return includesAssetMap
}

// ToIncludesLocalizedAssetMap returns a map of Asset's from the Includes
func (col *Collection) ToIncludesLocalizedAssetMap() map[string]*Asset {
	var includesAsset []*Asset
	includesAssetMap := make(map[string]*Asset)

	byteArray, _ := json.Marshal(col.Includes["Asset"])
	json.NewDecoder(bytes.NewReader(byteArray)).Decode(&includesAsset)

	for _, a := range includesAsset {
		includesAssetMap[a.Sys.ID] = a
	}
	return includesAssetMap
}
