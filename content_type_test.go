package contentful

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleContentTypesService_Get() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(context.TODO(), "space-id", "content-type-id")
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	fmt.Println(contentType.Name)
}

func ExampleContentTypesService_List() {
	cma := NewCMA("cma-token")

	collection, err := cma.ContentTypes.List(context.TODO(), "space-id").Next()
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	contentTypes, err := collection.ToContentType()
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	for _, contentType := range contentTypes {
		fmt.Println(contentType.Sys.ID, contentType.Sys.PublishedAt)
	}
}

func ExampleContentTypesService_Upsert_create() {
	cma := NewCMA("cma-token")

	contentType := &ContentType{
		Name:         "test content type",
		DisplayField: "field1_id",
		Description:  "content type description",
		Fields: []*Field{
			{
				ID:       "field1_id",
				Name:     "field1",
				Type:     "Symbol",
				Required: false,
				Disabled: false,
			},
			{
				ID:       "field2_id",
				Name:     "field2",
				Type:     "Symbol",
				Required: false,
				Disabled: true,
			},
		},
	}

	err := cma.ContentTypes.Upsert(context.TODO(), "space-id", contentType)
	if err != nil {
		log.Fatal(err) //nolint:revive
	}
}

func ExampleContentTypesService_Upsert_update() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(context.TODO(), "space-id", "content-type-id")
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	contentType.Name = "modified content type name"

	err = cma.ContentTypes.Upsert(context.TODO(), "space-id", contentType)
	if err != nil {
		log.Fatal(err) //nolint:revive
	}
}

func ExampleContentTypesService_Activate() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(context.TODO(), "space-id", "content-type-id")
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	err = cma.ContentTypes.Activate(context.TODO(), "space-id", contentType)
	if err != nil {
		log.Fatal(err) //nolint:revive
	}
}

func ExampleContentTypesService_Deactivate() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(context.TODO(), "space-id", "content-type-id")
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	err = cma.ContentTypes.Deactivate(context.TODO(), "space-id", contentType)
	if err != nil {
		log.Fatal(err) //nolint:revive
	}
}

func ExampleContentTypesService_Delete() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(context.TODO(), "space-id", "content-type-id")
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	err = cma.ContentTypes.Delete(context.TODO(), "space-id", contentType)
	if err != nil {
		log.Fatal(err) //nolint:revive
	}
}

func ExampleContentTypesService_Delete_allDrafts() {
	cma := NewCMA("cma-token")

	collection, err := cma.ContentTypes.List(context.TODO(), "space-id").Next()
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	contentTypes, err := collection.ToContentType()
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	for _, contentType := range contentTypes {
		if contentType.Sys.PublishedAt == "" {
			err := cma.ContentTypes.Delete(context.TODO(), "space-id", contentType)
			if err != nil {
				log.Fatal(err) //nolint:revive
			}
		}
	}
}

func TestContentTypesServiceList(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types", r.URL.Path)

		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_types.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.ContentTypes.List(context.TODO(), spaceID).Next()
	require.NoError(t, err)
}

func TestContentTypesServiceActivate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types/63Vgs0BFK0USe4i2mQUGK6/published", r.URL.Path)

		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData(t, "content_type.json")
	require.NoError(t, err)

	err = cma.ContentTypes.Activate(context.TODO(), spaceID, ct)
	require.NoError(t, err)
}

func TestContentTypesServiceDeactivate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types/63Vgs0BFK0USe4i2mQUGK6/published", r.URL.Path)

		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData(t, "content_type.json")
	require.NoError(t, err)

	err = cma.ContentTypes.Deactivate(context.TODO(), spaceID, ct)
	require.NoError(t, err)
}

func TestContentTypeSaveForCreate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)
		assert.Equal(t, "ct-name", payload["name"])
		assert.Equal(t, "ct-description", payload["description"])

		fields := payload["fields"].([]interface{})
		assert.Len(t, fields, 2)

		field1 := fields[0].(map[string]interface{})
		field2 := fields[1].(map[string]interface{})

		assert.Equal(t, "field1", field1["id"].(string))
		assert.Equal(t, "field1-name", field1["name"].(string))
		assert.Equal(t, "Symbol", field1["type"].(string))

		assert.Equal(t, "field2", field2["id"].(string))
		assert.Equal(t, "field2-name", field2["name"].(string))
		assert.Equal(t, "Symbol", field2["type"].(string))
		assert.True(t, field2["disabled"].(bool))

		assert.Equal(t, field1["id"].(string), payload["displayField"])

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:       "field1",
		Name:     "field1-name",
		Type:     "Symbol",
		Required: true,
	}

	field2 := &Field{
		ID:       "field2",
		Name:     "field2-name",
		Type:     "Symbol",
		Disabled: true,
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1, field2},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
	assert.Equal(t, "63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
	assert.Equal(t, "ct-name", ct.Name)
	assert.Equal(t, "ct-description", ct.Description)
}

func TestContentTypeSaveForUpdate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types/63Vgs0BFK0USe4i2mQUGK6", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)
		assert.Equal(t, "ct-name-updated", payload["name"])
		assert.Equal(t, "ct-description-updated", payload["description"])

		fields := payload["fields"].([]interface{})
		assert.Len(t, fields, 3)

		field1 := fields[0].(map[string]interface{})
		field2 := fields[1].(map[string]interface{})
		field3 := fields[2].(map[string]interface{})

		assert.Equal(t, "field1", field1["id"].(string))
		assert.Equal(t, "field1-name-updated", field1["name"].(string))
		assert.Equal(t, "String", field1["type"].(string))

		assert.Equal(t, "field2", field2["id"].(string))
		assert.Equal(t, "field2-name-updated", field2["name"].(string))
		assert.Equal(t, "Integer", field2["type"].(string))
		assert.Nil(t, field2["disabled"])

		assert.Equal(t, "field3", field3["id"].(string))
		assert.Equal(t, "field3-name", field3["name"].(string))
		assert.Equal(t, "Date", field3["type"].(string))

		assert.Equal(t, field3["id"].(string), payload["displayField"])

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData(t, "content_type.json")
	require.NoError(t, err)

	ct.Name = "ct-name-updated"
	ct.Description = "ct-description-updated"

	field1 := ct.Fields[0]
	field1.Name = "field1-name-updated"
	field1.Type = "String"
	field1.Required = false

	field2 := ct.Fields[1]
	field2.Name = "field2-name-updated"
	field2.Type = "Integer"
	field2.Disabled = false

	field3 := &Field{
		ID:   "field3",
		Name: "field3-name",
		Type: "Date",
	}

	ct.Fields = append(ct.Fields, field3)
	ct.DisplayField = ct.Fields[2].ID

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
	assert.Equal(t, "63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
	assert.Equal(t, "ct-name-updated", ct.Name)
	assert.Equal(t, "ct-description-updated", ct.Description)
	assert.Equal(t, 2, ct.Sys.Version)
}

func TestContentTypeCreateWithoutID(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/id1/content_types", r.RequestURI)
		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct := &ContentType{
		Sys:  &Sys{},
		Name: "MyContentType",
	}

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
}

func TestContentTypeCreateWithID(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/spaces/id1/content_types/mycontenttype", r.RequestURI)
		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct := &ContentType{
		Sys: &Sys{
			ID: "mycontenttype",
		},
		Name: "MyContentType",
	}

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
}

func TestContentTypeDelete(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types/63Vgs0BFK0USe4i2mQUGK6", r.RequestURI)
		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData(t, "content_type.json")
	require.NoError(t, err)

	// delete content type
	err = cma.ContentTypes.Delete(context.TODO(), "id1", ct)
	require.NoError(t, err)
}

func TestContentTypeFieldRef(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		fields := payload["fields"].([]interface{})
		assert.Len(t, fields, 1)

		field1 := fields[0].(map[string]interface{})
		assert.Equal(t, "Link", field1["type"].(string))
		validations := field1["validations"].([]interface{})
		assert.Len(t, validations, 1)
		validation := validations[0].(map[string]interface{})
		linkValidationValue := validation["linkContentType"].([]interface{})
		assert.Len(t, linkValidationValue, 1)
		assert.Equal(t, "63Vgs0BFK0USe4i2mQUGK6", linkValidationValue[0].(string))

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	linkCt, err := contentTypeFromTestData(t, "content_type.json")
	require.NoError(t, err)

	field1 := &Field{
		ID:   "field1",
		Name: "field1-name",
		Type: "Link",
		Validations: []FieldValidation{
			FieldValidationLink{
				LinkContentType: []string{linkCt.Sys.ID},
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
}

func TestContentTypeFieldArray(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		fields := payload["fields"].([]interface{})
		assert.Len(t, fields, 1)

		field1 := fields[0].(map[string]interface{})
		assert.Equal(t, "Array", field1["type"].(string))

		arrayItemSchema := field1["items"].(map[string]interface{})
		assert.Equal(t, "Text", arrayItemSchema["type"].(string))

		arrayItemSchemaValidations := arrayItemSchema["validations"].([]interface{})
		validation1 := arrayItemSchemaValidations[0].(map[string]interface{})
		assert.True(t, validation1["unique"].(bool))

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:   "field1",
		Name: "field1-name",
		Type: FieldTypeArray,
		Items: &FieldTypeArrayItem{
			Type: FieldTypeText,
			Validations: []FieldValidation{
				&FieldValidationUnique{
					Unique: true,
				},
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
}

func TestContentTypeFieldValidationRangeUniquePredefinedValues(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		fields := payload["fields"].([]interface{})
		assert.Len(t, fields, 1)

		field1 := fields[0].(map[string]interface{})
		assert.Equal(t, "Integer", field1["type"].(string))

		validations := field1["validations"].([]interface{})

		// unique validation
		validationUnique := validations[0].(map[string]interface{})
		assert.False(t, validationUnique["unique"].(bool))

		// range validation
		validationRange := validations[1].(map[string]interface{})
		rangeValues := validationRange["range"].(map[string]interface{})
		errorMessage := validationRange["message"].(string)
		assert.Equal(t, "error message", errorMessage)
		assert.Equal(t, float64(20), rangeValues["min"].(float64))
		assert.Equal(t, float64(30), rangeValues["max"].(float64))

		// predefined validation
		validationPredefinedValues := validations[2].(map[string]interface{})
		predefinedValues := validationPredefinedValues["in"].([]interface{})
		assert.Len(t, predefinedValues, 3)
		assert.Equal(t, "error message 2", validationPredefinedValues["message"].(string))
		assert.Equal(t, float64(20), predefinedValues[0].(float64))
		assert.Equal(t, float64(21), predefinedValues[1].(float64))
		assert.Equal(t, float64(22), predefinedValues[2].(float64))

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:   "field1",
		Name: "field1-name",
		Type: FieldTypeInteger,
		Validations: []FieldValidation{
			&FieldValidationUnique{
				Unique: false,
			},
			&FieldValidationRange{
				Range: &MinMax{
					Min: 20,
					Max: 30,
				},
				ErrorMessage: "error message",
			},
			&FieldValidationPredefinedValues{
				In:           []interface{}{20, 21, 22},
				ErrorMessage: "error message 2",
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
}

func TestContentTypeFieldTypeMedia(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/content_types", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		fields := payload["fields"].([]interface{})
		assert.Len(t, fields, 1)

		field1 := fields[0].(map[string]interface{})
		assert.Equal(t, "Link", field1["type"].(string))
		assert.Equal(t, "Asset", field1["linkType"].(string))

		validations := field1["validations"].([]interface{})

		// mime type validation
		validationMimeType := validations[0].(map[string]interface{})
		linkMimetypeGroup := validationMimeType["linkMimetypeGroup"].([]interface{})
		assert.Len(t, linkMimetypeGroup, 12)
		var mimetypes []string
		for _, mimetype := range linkMimetypeGroup {
			mimetypes = append(mimetypes, mimetype.(string))
		}
		assert.Equal(t, []string{
			MimeTypeAttachment,
			MimeTypePlainText,
			MimeTypeImage,
			MimeTypeAudio,
			MimeTypeVideo,
			MimeTypeRichText,
			MimeTypePresentation,
			MimeTypeSpreadSheet,
			MimeTypePDF,
			MimeTypeArchive,
			MimeTypeCode,
			MimeTypeMarkup,
		}, mimetypes)

		// dimension validation
		validationDimension := validations[1].(map[string]interface{})
		errorMessage := validationDimension["message"].(string)
		assetImageDimensions := validationDimension["assetImageDimensions"].(map[string]interface{})
		widthData := assetImageDimensions["width"].(map[string]interface{})
		heightData := assetImageDimensions["height"].(map[string]interface{})
		widthMin := int(widthData["min"].(float64))
		heightMax := int(heightData["max"].(float64))

		_, ok := widthData["max"].(float64)
		assert.False(t, ok)

		_, ok = heightData["min"].(float64)
		assert.False(t, ok)

		assert.Equal(t, "custom error message", errorMessage)
		assert.Equal(t, 100, widthMin)
		assert.Equal(t, 300, heightMax)

		// size validation
		validationSize := validations[2].(map[string]interface{})
		sizeData := validationSize["assetFileSize"].(map[string]interface{})
		minimum := int(sizeData["min"].(float64))
		maximum := int(sizeData["max"].(float64))
		assert.Equal(t, 30, minimum)
		assert.Equal(t, 400, maximum)

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:       "field-id",
		Name:     "media-field",
		Type:     FieldTypeLink,
		LinkType: "Asset",
		Validations: []FieldValidation{
			&FieldValidationMimeType{
				MimeTypes: []string{
					MimeTypeAttachment,
					MimeTypePlainText,
					MimeTypeImage,
					MimeTypeAudio,
					MimeTypeVideo,
					MimeTypeRichText,
					MimeTypePresentation,
					MimeTypeSpreadSheet,
					MimeTypePDF,
					MimeTypeArchive,
					MimeTypeCode,
					MimeTypeMarkup,
				},
			},
			&FieldValidationDimension{
				Width: &MinMax{
					Min: 100,
				},
				Height: &MinMax{
					Max: 300,
				},
				ErrorMessage: "custom error message",
			},
			&FieldValidationFileSize{
				Size: &MinMax{
					Min: 30,
					Max: 400,
				},
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(context.TODO(), "id1", ct)
	require.NoError(t, err)
}

func TestContentTypeFieldValidationsUnmarshal(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "content_type_with_validations.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	ct, err := cma.ContentTypes.Get(context.TODO(), spaceID, "validationsTest")
	require.NoError(t, err)

	uniqueValidations := []FieldValidation{}
	linkValidations := []FieldValidation{}
	sizeValidations := []FieldValidation{}
	regexValidations := []FieldValidation{}
	preDefinedValidations := []FieldValidation{}
	rangeValidations := []FieldValidation{}
	dateValidations := []FieldValidation{}
	mimeTypeValidations := []FieldValidation{}
	dimensionValidations := []FieldValidation{}
	fileSizeValidations := []FieldValidation{}

	for _, field := range ct.Fields {
		if field.Name == "text-short" {
			assert.Len(t, field.Validations, 4)
			uniqueValidations = append(uniqueValidations, field.Validations[0])
			sizeValidations = append(sizeValidations, field.Validations[1])
			regexValidations = append(regexValidations, field.Validations[2])
			preDefinedValidations = append(preDefinedValidations, field.Validations[3])
		}

		if field.Name == "text-long" {
			assert.Len(t, field.Validations, 3)
			sizeValidations = append(sizeValidations, field.Validations[0])
			regexValidations = append(regexValidations, field.Validations[1])
			preDefinedValidations = append(preDefinedValidations, field.Validations[2])
		}

		if field.Name == "number-integer" || field.Name == "number-decimal" {
			assert.Len(t, field.Validations, 3)
			uniqueValidations = append(uniqueValidations, field.Validations[0])
			rangeValidations = append(rangeValidations, field.Validations[1])
			preDefinedValidations = append(preDefinedValidations, field.Validations[2])
		}

		if field.Name == "date" {
			assert.Len(t, field.Validations, 1)
			dateValidations = append(dateValidations, field.Validations[0])
		}

		if field.Name == "location" || field.Name == "bool" {
			assert.Empty(t, field.Validations)
		}

		if field.Name == "media-onefile" {
			assert.Len(t, field.Validations, 3)
			mimeTypeValidations = append(mimeTypeValidations, field.Validations[0])
			dimensionValidations = append(dimensionValidations, field.Validations[1])
			fileSizeValidations = append(fileSizeValidations, field.Validations[2])
		}

		if field.Name == "media-manyfiles" {
			assert.Len(t, field.Validations, 1)
			assert.Len(t, field.Items.Validations, 3)
			sizeValidations = append(sizeValidations, field.Validations[0])
			mimeTypeValidations = append(mimeTypeValidations, field.Items.Validations[0])
			dimensionValidations = append(dimensionValidations, field.Items.Validations[1])
			fileSizeValidations = append(fileSizeValidations, field.Items.Validations[2])
		}

		if field.Name == "json" {
			assert.Len(t, field.Validations, 1)
			sizeValidations = append(sizeValidations, field.Validations[0])
		}

		if field.Name == "ref-onref" {
			assert.Len(t, field.Validations, 1)
			linkValidations = append(linkValidations, field.Validations[0])
		}

		if field.Name == "ref-manyRefs" {
			assert.Len(t, field.Validations, 1)
			assert.Len(t, field.Items.Validations, 1)
			linkValidations = append(linkValidations, field.Items.Validations[0])
			sizeValidations = append(sizeValidations, field.Validations[0])
		}
	}

	for _, validation := range uniqueValidations {
		_, ok := validation.(FieldValidationUnique)
		assert.True(t, ok)
	}

	for _, validation := range linkValidations {
		_, ok := validation.(FieldValidationLink)
		assert.True(t, ok)
	}

	for _, validation := range sizeValidations {
		_, ok := validation.(FieldValidationSize)
		assert.True(t, ok)
	}

	for _, validation := range regexValidations {
		_, ok := validation.(FieldValidationRegex)
		assert.True(t, ok)
	}

	for _, validation := range preDefinedValidations {
		_, ok := validation.(FieldValidationPredefinedValues)
		assert.True(t, ok)
	}

	for _, validation := range rangeValidations {
		_, ok := validation.(FieldValidationRange)
		assert.True(t, ok)
	}

	for _, validation := range dateValidations {
		_, ok := validation.(FieldValidationDate)
		assert.True(t, ok)
	}

	for _, validation := range mimeTypeValidations {
		_, ok := validation.(FieldValidationMimeType)
		assert.True(t, ok)
	}

	for _, validation := range dimensionValidations {
		_, ok := validation.(FieldValidationDimension)
		assert.True(t, ok)
	}

	for _, validation := range fileSizeValidations {
		_, ok := validation.(FieldValidationFileSize)
		assert.True(t, ok)
	}
}
