package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Sys    *Sys        `json:"sys"`
	Fields *FileFields `json:"fields"`
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
func (service *AssetsService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/assets", spaceID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single asset entity
func (service *AssetsService) Get(spaceID, assetID string, locale ...string) (*Asset, error) {
	path := fmt.Sprintf("/spaces/%s/assets/%s", spaceID, assetID)
	query := url.Values{}
	if service.c.api == "CDA" && len(locale) > 0 {
		query["locale"] = locale
	}

	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
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
		fmt.Println("LOCALE IS: ", retLocale)
		asset.Sys = assetNoLocale.Sys
		// asset.Fields.Title[retLocale] = assetNoLocale.Fields.Title
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

// Upsert updates or creates a new asset entity
func (service *AssetsService) Upsert(spaceID string, asset *Asset) error {
	bytesArray, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	var path string
	var method string

	if asset.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s/assets/%s", spaceID, asset.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/assets", spaceID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	return service.c.do(req, asset)
}

// Delete sends delete request
func (service *AssetsService) Delete(spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s/assets/%s", spaceID, asset.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Process the asset
func (service *AssetsService) Process(spaceID string, asset *Asset) error {
	var locale string
	for k := range asset.Fields.Title {
		locale = k
		break
	}

	path := fmt.Sprintf("/spaces/%s/assets/%s/files/%s/process", spaceID, asset.Sys.ID, locale)
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Publish publishes the asset
func (service *AssetsService) Publish(spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s/assets/%s/published", spaceID, asset.Sys.ID)
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, asset)
}

// Unpublish unpublishes the asset
func (service *AssetsService) Unpublish(spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s/assets/%s/published", spaceID, asset.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, asset)
}
