package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// WebhooksService service
type WebhooksService service

// Webhook model
type Webhook struct {
	Sys               *Sys             `json:"sys,omitempty"`
	Name              string           `json:"name,omitempty"`
	URL               string           `json:"url,omitempty"`
	Topics            []string         `json:"topics,omitempty"`
	HTTPBasicUsername string           `json:"httpBasicUsername,omitempty"`
	HTTPBasicPassword string           `json:"httpBasicPassword,omitempty"`
	Headers           []*WebhookHeader `json:"headers,omitempty"`
}

// WebhookHeader model
type WebhookHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetVersion returns entity version
func (webhook *Webhook) GetVersion() int {
	version := 1
	if webhook.Sys != nil {
		version = webhook.Sys.Version
	}

	return version
}

// List returns webhooks collection
func (service *WebhooksService) List(ctx context.Context, spaceID string) *Collection[Webhook] {
	path := fmt.Sprintf("/spaces/%s%s/webhook_definitions", spaceID, getEnvPath(service.c))
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return &Collection[Webhook]{}
	}

	col := NewCollection[Webhook](&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single webhook entity
func (service *WebhooksService) Get(ctx context.Context, spaceID, webhookID string) (*Webhook, error) {
	path := fmt.Sprintf("/spaces/%s%s/webhook_definitions/%s", spaceID, getEnvPath(service.c), webhookID)
	method := http.MethodGet

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	if err := service.c.do(req, &webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

// Upsert updates or creates a new entity
func (service *WebhooksService) Upsert(ctx context.Context, spaceID string, webhook *Webhook) error {
	bytesArray, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	var path string
	var method string

	if webhook.Sys != nil && webhook.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s%s/webhook_definitions/%s", spaceID, getEnvPath(service.c), webhook.Sys.ID)
		method = http.MethodPut
	} else {
		path = fmt.Sprintf("/spaces/%s%s/webhook_definitions", spaceID, getEnvPath(service.c))
		method = http.MethodPost
	}

	req, err := service.c.newRequest(ctx, method, path, nil, bytes.NewReader(bytesArray), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(webhook.GetVersion()))

	return service.c.do(req, webhook)
}

// Delete the webhook
func (service *WebhooksService) Delete(ctx context.Context, spaceID string, webhook *Webhook) error {
	path := fmt.Sprintf("/spaces/%s%s/webhook_definitions/%s", spaceID, getEnvPath(service.c), webhook.Sys.ID)
	method := http.MethodDelete

	req, err := service.c.newRequest(ctx, method, path, nil, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(webhook.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}
