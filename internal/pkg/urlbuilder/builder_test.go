package urlbuilder

import (
	"net/url"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

func TestBuilder_Build_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenNodes   []*html.Node
		givenURL     *url.URL
		expectedURLs []*url.URL
	}{
		{
			name: "given nodes, expect them to be parsed to URLs",
			givenNodes: []*html.Node{
				{
					Attr: []html.Attribute{
						{
							Key: "href",
							Val: "https://example.com/test",
						},
					},
				},
				{
					Attr: []html.Attribute{
						{
							Key: "href",
							Val: "/signup",
						},
					},
				},
				{
					Attr: []html.Attribute{
						{
							Key: "href",
							Val: "https://example.com/test?name=jamie",
						},
					},
				},
			},
			givenURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			expectedURLs: []*url.URL{
				{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/signup/",
				},
				{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/test/",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := New()

			actual, err := builder.Build(test.givenNodes, test.givenURL)
			if err != nil {
				t.Fatal(err)
			}

			sort.Slice(actual, func(i, j int) bool {
				return actual[i].Path < actual[j].Path
			})

			if !cmp.Equal(actual, test.expectedURLs) {
				t.Fatal(cmp.Diff(actual, test.expectedURLs))
			}
		})
	}
}

func TestBuilder_Build_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenNodes    []*html.Node
		givenURL      *url.URL
		expectedError error
	}{
		{
			name: "given an ASCII control character, expect an error",
			givenNodes: []*html.Node{
				{
					Attr: []html.Attribute{
						{
							Key: "href",
							Val: string(rune(0x7f)),
						},
					},
				},
				{
					Attr: []html.Attribute{
						{
							Key: "href",
							Val: "/signup",
						},
					},
				},
			},
			givenURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			expectedError: ErrFailedToParseURL,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := New()

			_, err := builder.Build(test.givenNodes, test.givenURL)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}
