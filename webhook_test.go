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

func TestWebhookSaveForCreate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/webhook_definitions", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "webhook-name", payload["name"])
		assert.Equal(t, "https://www.example.com/test", payload["url"])
		assert.Equal(t, "username", payload["httpBasicUsername"])
		assert.Equal(t, "password", payload["httpBasicPassword"])

		topics := payload["topics"].([]interface{})
		assert.Len(t, topics, 2)
		assert.Equal(t, "Entry.create", topics[0].(string))
		assert.Equal(t, "ContentType.create", topics[1].(string))

		headers := payload["headers"].([]interface{})
		assert.Len(t, headers, 2)
		header1 := headers[0].(map[string]interface{})
		header2 := headers[1].(map[string]interface{})

		assert.Equal(t, "header1", header1["key"].(string))
		assert.Equal(t, "header1-value", header1["value"].(string))

		assert.Equal(t, "header2", header2["key"].(string))
		assert.Equal(t, "header2-value", header2["value"].(string))

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, readTestData(t, "webhook.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	webhook := &Webhook{
		Name: "webhook-name",
		URL:  "https://www.example.com/test",
		Topics: []string{
			"Entry.create",
			"ContentType.create",
		},
		HTTPBasicUsername: "username",
		HTTPBasicPassword: "password",
		Headers: []*WebhookHeader{
			{
				Key:   "header1",
				Value: "header1-value",
			},
			{
				Key:   "header2",
				Value: "header2-value",
			},
		},
	}

	err = cma.Webhooks.Upsert(context.TODO(), spaceID, webhook)
	require.NoError(t, err)
	assert.Equal(t, "7fstd9fZ9T2p3kwD49FxhI", webhook.Sys.ID)
	assert.Equal(t, "webhook-name", webhook.Name)
	assert.Equal(t, "username", webhook.HTTPBasicUsername)
}

func TestWebhookSaveForUpdate(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/webhook_definitions/7fstd9fZ9T2p3kwD49FxhI", r.RequestURI)
		checkHeaders(t, r)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, "updated-webhook-name", payload["name"])
		assert.Equal(t, "https://www.example.com/test-updated", payload["url"])
		assert.Equal(t, "updated-username", payload["httpBasicUsername"])
		assert.Equal(t, "updated-password", payload["httpBasicPassword"])

		topics := payload["topics"].([]interface{})
		assert.Len(t, topics, 3)
		assert.Equal(t, "Entry.create", topics[0].(string))
		assert.Equal(t, "ContentType.create", topics[1].(string))
		assert.Equal(t, "Asset.create", topics[2].(string))

		headers := payload["headers"].([]interface{})
		assert.Len(t, headers, 2)
		header1 := headers[0].(map[string]interface{})
		header2 := headers[1].(map[string]interface{})

		assert.Equal(t, "header1", header1["key"].(string))
		assert.Equal(t, "updated-header1-value", header1["value"].(string))

		assert.Equal(t, "header2", header2["key"].(string))
		assert.Equal(t, "updated-header2-value", header2["value"].(string))

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, readTestData(t, "webhook-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test webhook
	webhook, err := webhookFromTestData(t, "webhook.json")
	require.NoError(t, err)

	webhook.Name = "updated-webhook-name"
	webhook.URL = "https://www.example.com/test-updated"
	webhook.Topics = []string{
		"Entry.create",
		"ContentType.create",
		"Asset.create",
	}
	webhook.HTTPBasicUsername = "updated-username"
	webhook.HTTPBasicPassword = "updated-password"
	webhook.Headers = []*WebhookHeader{
		{
			Key:   "header1",
			Value: "updated-header1-value",
		},
		{
			Key:   "header2",
			Value: "updated-header2-value",
		},
	}

	err = cma.Webhooks.Upsert(context.TODO(), spaceID, webhook)
	require.NoError(t, err)
	assert.Equal(t, "7fstd9fZ9T2p3kwD49FxhI", webhook.Sys.ID)
	assert.Equal(t, 1, webhook.Sys.Version)
	assert.Equal(t, "updated-webhook-name", webhook.Name)
	assert.Equal(t, "updated-username", webhook.HTTPBasicUsername)
}

func TestWebhookDelete(t *testing.T) {
	var err error

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/spaces/"+spaceID+"/webhook_definitions/7fstd9fZ9T2p3kwD49FxhI", r.RequestURI)
		checkHeaders(t, r)

		w.WriteHeader(http.StatusOK)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test webhook
	webhook, err := webhookFromTestData(t, "webhook.json")
	require.NoError(t, err)

	err = cma.Webhooks.Delete(context.TODO(), spaceID, webhook)
	require.NoError(t, err)
}
