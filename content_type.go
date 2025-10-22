package contentful

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// ContentTypesService service
type ContentTypesService service

// ContentType model
type ContentType struct {
	Sys          *Sys     `json:"sys"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Fields       []*Field `json:"fields,omitempty"`
	DisplayField string   `json:"displayField,omitempty"`
}

const (
	// FieldTypeSymbol content type field type for short textual data
	FieldTypeSymbol = "Symbol"

	// FieldTypeText content type field type for text data
	FieldTypeText = "Text"

	// FieldTypeArray content type field type for array data
	FieldTypeArray = "Array"

	// FieldTypeLink content type field type for link data
	FieldTypeLink = "Link"

	// FieldTypeInteger content type field type for integer data
	FieldTypeInteger = "Integer"

	// FieldTypeLocation content type field type for location data
	FieldTypeLocation = "Location"

	// FieldTypeBoolean content type field type for boolean data
	FieldTypeBoolean = "Boolean"

	// FieldTypeDate content type field type for date data
	FieldTypeDate = "Date"

	// FieldTypeObject content type field type for object data
	FieldTypeObject = "Object"
)

// Field model
type Field struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	LinkType    string              `json:"linkType,omitempty"`
	Items       *FieldTypeArrayItem `json:"items,omitempty"`
	Required    bool                `json:"required,omitempty"`
	Localized   bool                `json:"localized,omitempty"`
	Disabled    bool                `json:"disabled,omitempty"`
	Omitted     bool                `json:"omitted,omitempty"`
	Validations []FieldValidation   `json:"validations,omitempty"`
}

// UnmarshalJSON for custom json unmarshaling
func (field *Field) UnmarshalJSON(data []byte) error {
	payload := map[string]interface{}{}
	if err := Unmarshal(data, &payload); err != nil {
		return err
	}

	if val, ok := payload["id"]; ok {
		field.ID = val.(string)
	}

	if val, ok := payload["name"]; ok {
		field.Name = val.(string)
	}

	if val, ok := payload["type"]; ok {
		field.Type = val.(string)
	}

	if val, ok := payload["linkType"]; ok {
		field.LinkType = val.(string)
	}

	if val, ok := payload["items"]; ok {
		var fieldTypeArrayItem FieldTypeArrayItem
		if err := DeepCopy(&fieldTypeArrayItem, val); err != nil {
			return err
		}

		field.Items = &fieldTypeArrayItem
	}

	if val, ok := payload["required"]; ok {
		field.Required = val.(bool)
	}

	if val, ok := payload["localized"]; ok {
		field.Localized = val.(bool)
	}

	if val, ok := payload["disabled"]; ok {
		field.Disabled = val.(bool)
	}

	if val, ok := payload["omitted"]; ok {
		field.Omitted = val.(bool)
	}

	if val, ok := payload["validations"]; ok {
		validations, err := ParseValidations(val.([]interface{}))
		if err != nil {
			return err
		}

		field.Validations = validations
	}

	return nil
}

// ParseValidations converts json representation to go struct
func ParseValidations(data []interface{}) ([]FieldValidation, error) {
	var validations []FieldValidation
	for _, value := range data {
		buf, done := Buffer()
		var validation map[string]interface{}

		if validationStr, ok := value.(string); ok {
			buf.WriteString(validationStr)
			if err := Decode(buf, &validation); err != nil {
				done()
				return nil, err
			}
		}

		if validationMap, ok := value.(map[string]interface{}); ok {
			buf.Reset()
			if err := Encode(buf, validationMap); err != nil {
				done()
				return nil, err
			}
			validation = validationMap
		}

		if _, ok := validation["linkContentType"]; ok {
			var fieldValidationLink FieldValidationLink
			if err := Decode(buf, &fieldValidationLink); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationLink)
		}

		if _, ok := validation["linkMimetypeGroup"]; ok {
			var fieldValidationMimeType FieldValidationMimeType
			if err := Decode(buf, &fieldValidationMimeType); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationMimeType)
		}

		if _, ok := validation["assetImageDimensions"]; ok {
			var fieldValidationDimension FieldValidationDimension
			if err := Decode(buf, &fieldValidationDimension); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationDimension)
		}

		if _, ok := validation["assetFileSize"]; ok {
			var fieldValidationFileSize FieldValidationFileSize
			if err := Decode(buf, &fieldValidationFileSize); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationFileSize)
		}

		if _, ok := validation["unique"]; ok {
			var fieldValidationUnique FieldValidationUnique
			if err := Decode(buf, &fieldValidationUnique); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationUnique)
		}

		if _, ok := validation["in"]; ok {
			var fieldValidationPredefinedValues FieldValidationPredefinedValues
			if err := Decode(buf, &fieldValidationPredefinedValues); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationPredefinedValues)
		}

		if _, ok := validation["range"]; ok {
			var fieldValidationRange FieldValidationRange
			if err := Decode(buf, &fieldValidationRange); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationRange)
		}

		if _, ok := validation["dateRange"]; ok {
			var fieldValidationDate FieldValidationDate
			if err := Decode(buf, &fieldValidationDate); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationDate)
		}

		if _, ok := validation["size"]; ok {
			var fieldValidationSize FieldValidationSize
			if err := Decode(buf, &fieldValidationSize); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationSize)
		}

		if _, ok := validation["regexp"]; ok {
			var fieldValidationRegex FieldValidationRegex
			if err := Decode(buf, &fieldValidationRegex); err != nil {
				done()
				return nil, err
			}

			validations = append(validations, fieldValidationRegex)
		}
		done()
	}

	return validations, nil
}

// FieldTypeArrayItem model
type FieldTypeArrayItem struct {
	Type        string            `json:"type,omitempty"`
	Validations []FieldValidation `json:"validations,omitempty"`
	LinkType    string            `json:"linkType,omitempty"`
}

// UnmarshalJSON for custom json unmarshaling
func (item *FieldTypeArrayItem) UnmarshalJSON(data []byte) error {
	payload := map[string]interface{}{}
	if err := Unmarshal(data, &payload); err != nil {
		return err
	}

	if val, ok := payload["type"]; ok {
		item.Type = val.(string)
	}

	if val, ok := payload["validations"]; ok {
		validations, err := ParseValidations(val.([]interface{}))
		if err != nil {
			return err
		}

		item.Validations = validations
	}

	if val, ok := payload["linktype"]; ok {
		item.LinkType = val.(string)
	} else if val, ok := payload["linkType"]; ok {
		item.LinkType = val.(string)
	}

	return nil
}

// GetVersion returns entity version
func (ct *ContentType) GetVersion() int {
	version := 1
	if ct.Sys != nil {
		version = ct.Sys.Version
	}

	return version
}

// List return a content type collection
func (service *ContentTypesService) List(ctx context.Context, spaceID string) *Collection[ContentType] {
	path := fmt.Sprintf("/spaces/%s%s/content_types", spaceID, getEnvPath(service.c))

	req, err := service.c.newRequest(ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil
	}

	col := NewCollection[ContentType](&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get fetched a content type specified by `contentTypeID`
func (service *ContentTypesService) Get(ctx context.Context, spaceID, contentTypeID string) (*ContentType, error) {
	path := fmt.Sprintf("/spaces/%s%s/content_types/%s", spaceID, getEnvPath(service.c), contentTypeID)

	req, err := service.c.newRequest(ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var ct *ContentType
	if err = service.c.do(req, &ct); err != nil {
		return nil, err
	}

	return ct, nil
}

// Upsert updates or creates a new content type
func (service *ContentTypesService) Upsert(ctx context.Context, spaceID string, ct *ContentType) error {
	bytesArray, err := Marshal(ct)
	if err != nil {
		return err
	}

	var path string
	var method string

	if ct.Sys != nil && ct.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s%s/content_types/%s", spaceID, getEnvPath(service.c), ct.Sys.ID)
		method = http.MethodPut
	} else {
		path = fmt.Sprintf("/spaces/%s%s/content_types", spaceID, getEnvPath(service.c))
		method = http.MethodPost
	}

	req, err := service.c.newRequest(ctx, method, path, nil, bytes.NewReader(bytesArray), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(ct.GetVersion()))

	return service.c.do(req, &ct)
}

// Delete the content_type
func (service *ContentTypesService) Delete(ctx context.Context, spaceID string, ct *ContentType) error {
	path := fmt.Sprintf("/spaces/%s%s/content_types/%s", spaceID, getEnvPath(service.c), ct.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Activate the contenttype, a.k.a publish
func (service *ContentTypesService) Activate(ctx context.Context, spaceID string, ct *ContentType) error {
	path := fmt.Sprintf("/spaces/%s%s/content_types/%s/published", spaceID, getEnvPath(service.c), ct.Sys.ID)

	req, err := service.c.newRequest(ctx, http.MethodPut, path, nil, nil, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(ct.GetVersion()))

	return service.c.do(req, &ct)
}

// Deactivate the contenttype, a.k.a unpublish
func (service *ContentTypesService) Deactivate(ctx context.Context, spaceID string, ct *ContentType) error {
	path := fmt.Sprintf("/spaces/%s%s/content_types/%s/published", spaceID, getEnvPath(service.c), ct.Sys.ID)

	req, err := service.c.newRequest(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(ct.GetVersion()))

	return service.c.do(req, &ct)
}
