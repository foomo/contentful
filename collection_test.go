package contentful

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncTokenRegex(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		token string
	}{
		{
			name:  "empty",
			url:   "/spaces/yadj1kx9rmg0/sync?access_token=fdb4e7a3102747a02ea69ebac5e282b9e44d28fb340f778a4f5e788625a61abe",
			token: "",
		},
		{
			name:  "a-zA-Z0-9",
			url:   "/spaces/yadj1kx9rmg0/sync?access_token=fdb4e7a3102747a02ea69ebac5e282b9e44d28fb340f778a4f5e788625a61abe&sync_token=w7Ese3kdwpMbMhhgw7QAUsKiw6bCi09CwpFYwpwywqVYw6DDh8OawrTDpWvCgMOhw6jCuAhxWX_CocOPwowhcsOzeEJSbcOvwrfDlCjDr8O1YzLDvi9FOTXCmsOqT8OFcHPDuFDCqyMMTsKNw7rDmsOqKcOnw7FCwpIfNMOcFMOxFnHCoCzDpAjCucOdwpwfw4YTK8Kpw6zCtDrChVQlNsO2ZybDnw",
			token: "w7Ese3kdwpMbMhhgw7QAUsKiw6bCi09CwpFYwpwywqVYw6DDh8OawrTDpWvCgMOhw6jCuAhxWX_CocOPwowhcsOzeEJSbcOvwrfDlCjDr8O1YzLDvi9FOTXCmsOqT8OFcHPDuFDCqyMMTsKNw7rDmsOqKcOnw7FCwpIfNMOcFMOxFnHCoCzDpAjCucOdwpwfw4YTK8Kpw6zCtDrChVQlNsO2ZybDnw",
		},
		{
			name:  "a-zA-Z0-9_",
			url:   "/spaces/yadj1kx9rmg0/sync?access_token=fdb4e7a3102747a02ea69ebac5e282b9e44d28fb340f778a4f5e788625a61abe&sync_token=w_Ese3kdwpMbMhhgw7QAUsKiw6bCi09CwpFYwpwywqVYw6DDh8OawrTDpWvCgMOhw6jCuAhxWX_CocOPwowhcsOzeEJSbcOvwrfDlCjDr8O1YzLDvi9FOTXCmsOqT8OFcHPDuFDCqyMMTsKNw7rDmsOqKcOnw7FCwpIfNMOcFMOxFnHCoCzDpAjCucOdwpwfw4YTK8Kpw6zCtDrChVQlNsO2ZybDnw",
			token: "w_Ese3kdwpMbMhhgw7QAUsKiw6bCi09CwpFYwpwywqVYw6DDh8OawrTDpWvCgMOhw6jCuAhxWX_CocOPwowhcsOzeEJSbcOvwrfDlCjDr8O1YzLDvi9FOTXCmsOqT8OFcHPDuFDCqyMMTsKNw7rDmsOqKcOnw7FCwpIfNMOcFMOxFnHCoCzDpAjCucOdwpwfw4YTK8Kpw6zCtDrChVQlNsO2ZybDnw",
		},
		{
			name:  "a-zA-Z0-9_-",
			url:   "/spaces/yadj1kx9rmg0/sync?access_token=fdb4e7a3102747a02ea69ebac5e282b9e44d28fb340f778a4f5e788625a61abe&sync_token=w_-se3kdwpMbMhhgw7QAUsKiw6bCi09CwpFYwpwywqVYw6DDh8OawrTDpWvCgMOhw6jCuAhxWX_CocOPwowhcsOzeEJSbcOvwrfDlCjDr8O1YzLDvi9FOTXCmsOqT8OFcHPDuFDCqyMMTsKNw7rDmsOqKcOnw7FCwpIfNMOcFMOxFnHCoCzDpAjCucOdwpwfw4YTK8Kpw6zCtDrChVQlNsO2ZybDnw",
			token: "w_-se3kdwpMbMhhgw7QAUsKiw6bCi09CwpFYwpwywqVYw6DDh8OawrTDpWvCgMOhw6jCuAhxWX_CocOPwowhcsOzeEJSbcOvwrfDlCjDr8O1YzLDvi9FOTXCmsOqT8OFcHPDuFDCqyMMTsKNw7rDmsOqKcOnw7FCwpIfNMOcFMOxFnHCoCzDpAjCucOdwpwfw4YTK8Kpw6zCtDrChVQlNsO2ZybDnw",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matches := syncTokenRegex.FindStringSubmatch(test.url)
			if test.token == "" {
				assert.Empty(t, matches)
			} else {
				require.Len(t, matches, 2)
				assert.Equal(t, test.token, matches[1])
			}
		})
	}
}

func TestNewCollection(t *testing.T) {
	setup()
	defer teardown()
	t.Run("error notResolvable", func(t *testing.T) {
		c := NewCDA("asdasd")
		c.SetHTTPTransport(roundTrip{
			Filename: "error-notResolvable.json",
		})

		req, err := c.newRequest(context.TODO(), http.MethodGet, "/", url.Values{}, nil, nil)
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

		req, err := c.newRequest(context.TODO(), http.MethodGet, "/", url.Values{}, nil, nil)
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
