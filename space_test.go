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

func ExampleSpacesService_Get() {
	cma := NewCMA("cma-token")

	space, err := cma.Spaces.Get(context.TODO(), "space-id")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(space.Name)
}

func ExampleSpacesService_List() {
	cma := NewCMA("cma-token")
	collection, err := cma.Spaces.List(context.TODO()).Next()
	if err != nil {
		log.Fatal(err)
	}

	for _, space := range collection.Items {
		fmt.Println(space.Sys.ID, space.Name)
	}
}

func ExampleSpacesService_Upsert_create() {
	cma := NewCMA("cma-token")

	space := &Space{
		Name:          "space-name",
		DefaultLocale: "en-US",
	}

	err := cma.Spaces.Upsert(context.TODO(), space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Upsert_update() {
	cma := NewCMA("cma-token")

	space, err := cma.Spaces.Get(context.TODO(), "space-id")
	if err != nil {
		log.Fatal(err)
	}

	space.Name = "modified"
	err = cma.Spaces.Upsert(context.TODO(), space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Delete() {
	cma := NewCMA("cma-token")

	space, err := cma.Spaces.Get(context.TODO(), "space-id")
	if err != nil {
		log.Fatal(err)
	}

	err = cma.Spaces.Delete(context.TODO(), space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Delete_all() {
	cma := NewCMA("cma-token")

	collection, err := cma.Spaces.List(context.TODO()).Next()
	if err != nil {
		log.Fatal(err)
	}

	for _, space := range collection.Items {
		err := cma.Spaces.Delete(context.TODO(), space)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestSpacesServiceList(t *testing.T) {
	var err error
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/spaces", r.URL.Path)

		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "spaces.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Spaces.List(context.TODO()).Next()
	require.NoError(t, err)

	spaces := collection.Items
	require.NoError(t, err)
	assert.Len(t, spaces, 2)
	assert.Equal(t, "id1", spaces[0].Sys.ID)
	assert.Equal(t, "id2", spaces[1].Sys.ID)
}

func TestSpacesServiceGet(t *testing.T) {
	var err error
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/spaces/"+spaceID, r.URL.Path)

		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "space-1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space, err := cma.Spaces.Get(context.TODO(), spaceID)
	require.NoError(t, err)
	assert.Equal(t, "id1", space.Sys.ID)
}

func TestSpaceSaveForCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "new space", payload["name"])
		assert.Equal(t, "en", payload["defaultLocale"])

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "spaces-newspace.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space := &Space{
		Name:          "new space",
		DefaultLocale: "en",
	}

	err := cma.Spaces.Upsert(context.TODO(), space)
	require.NoError(t, err)
	assert.Equal(t, "newspace", space.Sys.ID)
	assert.Equal(t, "new space", space.Name)
	assert.Equal(t, "en", space.DefaultLocale)
}

func TestSpaceSaveForUpdate(t *testing.T) {
	var err error
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/spaces/newspace", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "changed-space-name", payload["name"])
		assert.Equal(t, "de", payload["defaultLocale"])

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "spaces-newspace-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space, err := spaceFromTestData(t, "spaces-newspace.json")
	require.NoError(t, err)

	space.Name = "changed-space-name"
	space.DefaultLocale = "de"

	err = cma.Spaces.Upsert(context.TODO(), space)
	require.NoError(t, err)
	assert.Equal(t, "changed-space-name", space.Name)
	assert.Equal(t, "de", space.DefaultLocale)
	assert.Equal(t, 2, space.Sys.Version)
}

func TestSpaceDelete(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/spaces/"+spaceID, r.RequestURI)
		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space, err := spaceFromTestData(t, "spaces-"+spaceID+".json")
	require.NoError(t, err)

	err = cma.Spaces.Delete(context.TODO(), space)
	require.NoError(t, err)
}
