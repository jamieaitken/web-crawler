package requester

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

func TestRequester_Build(t *testing.T) {
	tests := []struct {
		name            string
		givenURL        url.URL
		givenBody       io.Reader
		givenMethod     string
		expectedRequest *http.Request
	}{
		{
			name: "given a URL, method and body; expect them to be set within builder",
			givenURL: url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			givenBody:   strings.NewReader("example payload"),
			givenMethod: http.MethodPost,
			expectedRequest: &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
				},
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          io.NopCloser(strings.NewReader("example payload")),
				ContentLength: 15,
				Host:          "example.com",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := New().
				WithURL(test.givenURL).
				WithMethod(test.givenMethod).
				WithBody(test.givenBody).
				Build(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedRequest, options) {
				t.Fatal(cmp.Diff(actual, test.expectedRequest, options))
			}
		})
	}
}

var options = cmp.Options{
	cmpopts.IgnoreUnexported(strings.Reader{}, http.Request{}),
	cmpopts.IgnoreFields(http.Request{}, "GetBody"),
}
