package contentful

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestNewCollection(t *testing.T) {
	setup()
	defer teardown()
	t.Run("error notResolvable", func(t *testing.T) {
		c := NewCDA("asdasd")
		c.SetHTTPTransport(roundTrip{
			Filename: "error-notResolvable.json",
		})

		req, err := c.newRequest(context.TODO(), http.MethodGet, "/", url.Values{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		var data Collection
		if err := c.do(req, &data); err != nil {
			t.Fatal(err)
		}
		wantErrors := []Error{
			{
				Sys: &Sys{
					ID:   "notResolvable",
					Type: "error",
				},
				Details: map[string]string{
					"type":     "Link",
					"linkType": "Asset",
					"id":       "6J9XxxIv14p559rJwZjPgb",
				},
			},
		}
		if !reflect.DeepEqual(wantErrors, data.Errors) {
			t.Error("data.Errors has an incorrect structure.")
		}
	})
	t.Run("error nameUnknown", func(t *testing.T) {
		c := NewCDA("asdasd")
		c.SetHTTPTransport(roundTrip{
			Filename: "error-nameUnknown.json",
		})

		req, err := c.newRequest(context.TODO(), http.MethodGet, "/", url.Values{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		var data Collection
		if err := c.do(req, &data); err != nil {
			t.Fatal(err)
		}
		wantErrors := &ErrorDetails{
			Errors: []*ErrorDetail{
				{
					Name: "unknown",
					Path: []interface{}{
						"name",
					},
					Details: "The path \"name\" is not recognized",
				},
			},
		}

		if !reflect.DeepEqual(wantErrors, data.Details) {
			t.Error("data.Details has an incorrect structure.")
		}
	})
}

type roundTrip struct {
	Filename string      // required, must be a valid file in testdata directory
	Code     int         // optional, defaults to 200
	Header   http.Header // optional
}

func (rt roundTrip) RoundTrip(*http.Request) (*http.Response, error) {
	if rt.Code == 0 {
		rt.Code = http.StatusOK //nolint:revive
	}
	f, err := os.Open("testdata/" + rt.Filename)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		Header:     rt.Header,
		StatusCode: rt.Code,
		Body:       f,
	}, nil
}
