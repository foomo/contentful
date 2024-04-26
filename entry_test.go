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

func ExampleEntriesService_Upsert_create() {
	cma := NewCMA("cma-token")

	entry := &Entry{
		Sys: &Sys{
			ID: "MyEntry",
			ContentType: &ContentType{
				Sys: &Sys{
					ID: "MyContentType",
				},
			},
		},
		Fields: map[string]interface{}{
			"Description": map[string]string{
				"en-US": "Some example content...",
			},
		},
	}

	err := cma.Entries.Upsert(context.TODO(), "space-id", entry)
	if err != nil {
		log.Fatal(err) //nolint:revive
	}
}

func ExampleEntriesService_Upsert_update() {
	cma := NewCMA("cma-token")

	entry, err := cma.Entries.Get(context.TODO(), "space-id", "entry-id")
	if err != nil {
		log.Fatal(err) //nolint:revive
	}

	entry.Fields["Description"] = map[string]interface{}{
		"en-US": "modified entry content",
	}

	err = cma.Entries.Upsert(context.TODO(), "space-id", entry)
	if err != nil {
		log.Fatal(err) //nolint:revive
	}
}

func TestEntrySaveForCreate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/entries", r.RequestURI)
		assert.Equal(t, []string{"MyContentType"}, r.Header["X-Contentful-Content-Type"])
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		assert.NotNil(t, payload["fields"])
		fields := payload["fields"].(map[string]interface{})

		assert.Equal(t, map[string]interface{}{"en-US": "Some test content..."}, fields["Description"])

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "entry_3.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entry := &Entry{
		Sys: &Sys{
			ContentType: &ContentType{
				Sys: &Sys{
					ID: "MyContentType",
				},
			},
		},
		Fields: map[string]interface{}{
			"Description": map[string]string{
				"en-US": "Some test content...",
			},
		},
	}

	err = cma.Entries.Upsert(context.TODO(), "id1", entry)
	require.NoError(t, err)
	assert.Equal(t, "foocat", entry.Sys.ID)
}

// func TestEntrySaveForUpdate(t *testing.T) {
//	var err error
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, r.Method, "PUT")
//		assert.Equal(t, r.RequestURI, "/spaces/"+spaceID+"/content_types/63Vgs0BFK0USe4i2mQUGK6")
//		checkHeaders(t, r)
//
//		var payload map[string]interface{}
//		err := json.NewDecoder(r.Body).Decode(&payload)
//		require.NoError(t, err)
//		assert.Equal(t, "ct-name-updated", payload["name"])
//		assert.Equal(t, "ct-description-updated", payload["description"])
//
//		fields := payload["fields"].([]interface{})
//		assert.Equal(t, 3, len(fields))
//
//		field1 := fields[0].(map[string]interface{})
//		field2 := fields[1].(map[string]interface{})
//		field3 := fields[2].(map[string]interface{})
//
//		assert.Equal(t, "field1", field1["id"].(string))
//		assert.Equal(t, "field1-name-updated", field1["name"].(string))
//		assert.Equal(t, "String", field1["type"].(string))
//
//		assert.Equal(t, "field2", field2["id"].(string))
//		assert.Equal(t, "field2-name-updated", field2["name"].(string))
//		assert.Equal(t, "Integer", field2["type"].(string))
//		assert.Nil(field2["disabled"])
//
//		assert.Equal(t, "field3", field3["id"].(string))
//		assert.Equal(t, "field3-name", field3["name"].(string))
//		assert.Equal(t, "Date", field3["type"].(string))
//
//		assert.Equal(t, field3["id"].(string), payload["displayField"])
//
//		w.WriteHeader(200)
//		_, _ = fmt.Fprintln(w, readTestData(t, "content_type-updated.json"))
//	})
//
//	// test server
//	server := httptest.NewServer(handler)
//	defer server.Close()
//
//	// cma client
//	cma = NewCMA(CMAToken)
//	cma.BaseURL = server.URL
//
//	// test content type
//	ct, err := contentTypeFromTestData(t, "content_type.json")
//	require.NoError(t, err)
//
//	ct.Name = "ct-name-updated"
//	ct.Description = "ct-description-updated"
//
//	field1 := ct.Fields[0]
//	field1.Name = "field1-name-updated"
//	field1.Type = "String"
//	field1.Required = false
//
//	field2 := ct.Fields[1]
//	field2.Name = "field2-name-updated"
//	field2.Type = "Integer"
//	field2.Disabled = false
//
//	field3 := &Field{
//		ID:   "field3",
//		Name: "field3-name",
//		Type: "Date",
//	}
//
//	ct.Fields = append(ct.Fields, field3)
//	ct.DisplayField = ct.Fields[2].ID
//
//	cma.ContentTypes.Upsert("id1", ct)
//	require.NoError(t, err)
//	assert.Equal(t, "63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
//	assert.Equal(t, "ct-name-updated", ct.Name)
//	assert.Equal(t, "ct-description-updated", ct.Description)
//	assert.Equal(t, 2, ct.Sys.Version)
// }
//
// func TestEntryCreateWithoutID(t *testing.T) {
//	var err error
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, http.MethodPost, r.Method)
//		assert.Equal(t, "/spaces/id1/content_types", r.RequestURI)
//		checkHeaders(t, r)
//
//		w.WriteHeader(200)
//		_, _ = fmt.Fprintln(w, readTestData(t, "content_type-updated.json"))
//	})
//
//	// test server
//	server := httptest.NewServer(handler)
//	defer server.Close()
//
//	// cma client
//	cma = NewCMA(CMAToken)
//	cma.BaseURL = server.URL
//
//	// test content type
//	ct := &ContentType{
//		Sys:  &Sys{},
//		Name: "MyContentType",
//	}
//
//	cma.ContentTypes.Upsert("id1", ct)
//	require.NoError(t, err)
// }
//
// func TestEntryCreateWithID(t *testing.T) {
//	var err error
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, r.Method, "PUT")
//		assert.Equal(t, r.RequestURI, "/spaces/id1/content_types/mycontenttype")
//		checkHeaders(t, r)
//
//		w.WriteHeader(200)
//		_, _ = fmt.Fprintln(w, readTestData(t, "content_type-updated.json"))
//	})
//
//	// test server
//	server := httptest.NewServer(handler)
//	defer server.Close()
//
//	// cma client
//	cma = NewCMA(CMAToken)
//	cma.BaseURL = server.URL
//
//	// test content type
//	ct := &ContentType{
//		Sys: &Sys{
//			ID: "mycontenttype",
//		},
//		Name: "MyContentType",
//	}
//
//	cma.ContentTypes.Upsert("id1", ct)
//	require.NoError(t, err)
// }
