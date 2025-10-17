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

// AssetsService service
type AssetsService service

// File model
type File struct {
	Name        string      `json:"fileName,omitempty"`
	ContentType string      `json:"contentType,omitempty"`
	URL         string      `json:"url,omitempty"`
	UploadURL   string      `json:"upload,omitempty"`
	UploadFrom  *Upload     `json:"uploadFrom,omitempty"`
	Detail      *FileDetail `json:"details,omitempty"`
}

// FileDetail model
type FileDetail struct {
	Size  int        `json:"size,omitempty"`
	Image *FileImage `json:"image,omitempty"`
}

// FileImage model
type FileImage struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

// FileFields model
type FileFields struct {
	Title       map[string]string `json:"title,omitempty"`
	Description map[string]string `json:"description,omitempty"`
	File        map[string]*File  `json:"file,omitempty"`
}

// FileFieldsNoLocale model
type FileFieldsNoLocale struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	File        *File  `json:"file,omitempty"`
}

// Asset model
type Asset struct {
	Metadata *Metadata   `json:"metadata,omitempty"`
	Sys      *Sys        `json:"sys"`
	Fields   *FileFields `json:"fields"`
}

// AssetNoLocale model
type AssetNoLocale struct {
	Sys    *Sys                `json:"sys"`
	Fields *FileFieldsNoLocale `json:"fields"`
}

// GetVersion returns entity version
func (asset *Asset) GetVersion() int {
	version := 1
	if asset.Sys != nil {
		version = asset.Sys.Version
	}

	return version
}

// List returns asset collection
func (service *AssetsService) List(ctx context.Context, spaceID string) *Collection[*Asset] {
	path := fmt.Sprintf("/spaces/%s%s/assets", spaceID, getEnvPath(service.c))
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return &Collection[*Asset]{}
	}

	col := NewCollection[*Asset](&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single asset entity
func (service *AssetsService) Get(ctx context.Context, spaceID, assetID string, locale ...string) (*Asset, error) {
	path := fmt.Sprintf("/spaces/%s%s/assets/%s", spaceID, getEnvPath(service.c), assetID)
	query := url.Values{}
	if service.c.api == "CDA" && len(locale) > 0 {
		query["locale"] = locale
	}

	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	var asset Asset
	var assetNoLocale AssetNoLocale
	if service.c.api == "CDA" && (len(locale) == 0 || (len(locale) == 1 && locale[0] != "*")) {
		if err := service.c.do(req, &assetNoLocale); err != nil {
			return nil, err
		}
		retLocale := assetNoLocale.Sys.Locale
		asset.Sys = assetNoLocale.Sys
		localizedTitle := map[string]string{
			retLocale: assetNoLocale.Fields.Title,
		}
		localizedDescription := map[string]string{
			retLocale: assetNoLocale.Fields.Description,
		}
		localizedFile := map[string]*File{
			retLocale: assetNoLocale.Fields.File,
		}
		asset.Fields = &FileFields{
			Title:       localizedTitle,
			Description: localizedDescription,
			File:        localizedFile,
		}
		return &asset, nil
	}
	if err := service.c.do(req, &asset); err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetLocalized returns the asset with fields without localization map
func (asset *Asset) GetLocalized() *AssetNoLocale {
	// No default locale available if asking for all locales, fallback to nil
	if asset.Sys.Locale == "" {
		return nil
	}

	return &AssetNoLocale{
		Sys: asset.Sys,
		Fields: &FileFieldsNoLocale{
			Title:       asset.Fields.Title[asset.Sys.Locale],
			Description: asset.Fields.Description[asset.Sys.Locale],
			File:        asset.Fields.File[asset.Sys.Locale],
		},
	}
}

// Upsert updates or creates a new asset entity
func (service *AssetsService) Upsert(ctx context.Context, spaceID string, asset *Asset) error {
	bytesArray, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	var path string
	var method string

	if asset.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s%s/assets/%s", spaceID, getEnvPath(service.c), asset.Sys.ID)
		method = http.MethodPut
	} else {
		path = fmt.Sprintf("/spaces/%s%s/assets", spaceID, getEnvPath(service.c))
		method = http.MethodPost
	}

	req, err := service.c.newRequest(ctx, method, path, nil, bytes.NewReader(bytesArray), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	return service.c.do(req, asset)
}

// Delete sends delete request
func (service *AssetsService) Delete(ctx context.Context, spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s%s/assets/%s", spaceID, getEnvPath(service.c), asset.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Process the asset
func (service *AssetsService) Process(ctx context.Context, spaceID string, asset *Asset) error {
	var locale string
	for k := range asset.Fields.Title {
		locale = k
		path := fmt.Sprintf("/spaces/%s%s/assets/%s/files/%s/process", spaceID, getEnvPath(service.c), asset.Sys.ID, locale)
		method := http.MethodPut

		req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
		if err != nil {
			return err
		}

		version := strconv.Itoa(asset.Sys.Version)
		req.Header.Set("X-Contentful-Version", version)
		err = service.c.do(req, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// Publish publishes the asset
func (service *AssetsService) Publish(ctx context.Context, spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s%s/assets/%s/published", spaceID, getEnvPath(service.c), asset.Sys.ID)
	method := http.MethodPut

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, asset)
}

// Unpublish unpublishes the asset
func (service *AssetsService) Unpublish(ctx context.Context, spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s%s/assets/%s/published", spaceID, getEnvPath(service.c), asset.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, asset)
}
