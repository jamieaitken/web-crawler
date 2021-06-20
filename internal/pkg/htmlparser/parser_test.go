package htmlparser

import (
	"crawler/internal/pkg/urlbuilder"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"golang.org/x/net/html"

	"github.com/google/go-cmp/cmp"
)

func TestHTMLParser_FetchLinks_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenHTML       io.Reader
		givenURL        *url.URL
		givenURLBuilder URLBuilder
		expectedURLs    []*url.URL
	}{
		{
			name:      "given valid HTML, expect URLs returned",
			givenHTML: strings.NewReader("<html><a href='https://example.com/hi'>Hi</a></html>"),
			givenURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/welcome",
			},
			givenURLBuilder: mockURLBuilder{
				GivenURLs: []*url.URL{
					{
						Scheme: "https",
						Host:   "example.com",
						Path:   "/hi",
					},
				},
			},
			expectedURLs: []*url.URL{
				{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/hi",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := New(test.givenURLBuilder)

			urls, err := parser.FetchLinks(test.givenHTML, test.givenURL)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(urls, test.expectedURLs) {
				t.Fatal(cmp.Diff(urls, test.expectedURLs))
			}
		})
	}
}

func TestHTMLParser_FetchLinks_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenHTML       io.Reader
		givenURL        *url.URL
		givenURLBuilder URLBuilder
		expectedError   error
	}{
		{
			name:      "given an error from URL Builder, expect error returned",
			givenHTML: strings.NewReader("<html><a href='https://example.com/hi'>Hi</a></html>"),
			givenURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/welcome",
			},
			givenURLBuilder: mockURLBuilder{
				GivenError: urlbuilder.ErrFailedToParseURL,
			},
			expectedError: urlbuilder.ErrFailedToParseURL,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := New(test.givenURLBuilder)

			_, err := parser.FetchLinks(test.givenHTML, test.givenURL)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockURLBuilder struct {
	GivenURLs  []*url.URL
	GivenError error
}

func (m mockURLBuilder) Build(_ []*html.Node, _ *url.URL) ([]*url.URL, error) {
	return m.GivenURLs, m.GivenError
}
