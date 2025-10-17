package contentful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalesServiceList(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/locales", r.URL.Path)

		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "locales.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Locales.List(context.TODO(), spaceID).Next()
	require.NoError(t, err)
}

func TestLocalesServiceGet(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ", r.URL.Path)

		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "locale_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale, err := cma.Locales.Get(context.TODO(), spaceID, "4aGeQYgByqQFJtToAOh2JJ")
	require.NoError(t, err)
	assert.Equal(t, "U.S. English", locale.Name)
	assert.Equal(t, "en-US", locale.Code)
}

func TestLocalesServiceUpsertCreate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/locales", r.RequestURI)

		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "German (Austria)", payload["name"])
		assert.Equal(t, "de-AT", payload["code"])

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "locale_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale := &Locale{
		Name: "German (Austria)",
		Code: "de-AT",
	}

	err = cma.Locales.Upsert(context.TODO(), spaceID, locale)
	require.NoError(t, err)
}

func TestLocalesServiceUpsertUpdate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ", r.RequestURI)

		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "modified-name", payload["name"])
		assert.Equal(t, "modified-code", payload["code"])

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "locale_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale, err := localeFromTestData(t, "locale_1.json")
	require.NoError(t, err)

	locale.Name = "modified-name"
	locale.Code = "modified-code"

	err = cma.Locales.Upsert(context.TODO(), spaceID, locale)
	require.NoError(t, err)
}

func TestLocalesServiceDelete(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ", r.RequestURI)
		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test locale
	locale, err := localeFromTestData(t, "locale_1.json")
	require.NoError(t, err)

	// delete locale
	err = cma.Locales.Delete(context.TODO(), spaceID, locale)
	require.NoError(t, err)
}
