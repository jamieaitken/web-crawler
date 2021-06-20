package crawler

import (
	"context"
	"crawler/internal/domain"
	"crawler/internal/pkg/htmlparser"
	"crawler/internal/pkg/httpclient"
	"crawler/storage/memory"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

func TestController_Start_Success(t *testing.T) {
	tests := []struct {
		name          string
		givenBaseURL  *url.URL
		givenRepo     RepositoryProvider
		givenClient   ClientProvider
		givenParser   Parser
		expectedPages []domain.Page
	}{
		{
			name: "given one page, expect 3 links to be returned",
			givenBaseURL: &url.URL{
				Host: "example.com",
			},
			givenRepo: mockRepo{
				GivenGetError: memory.ErrInvalidKey,
			},
			givenClient: &mockClient{
				GivenFetchResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader(htmlBody)),
				},
				GivenFetchError: httpclient.ErrFailedToBuildRequest,
			},
			givenParser: mockParser{
				GivenURLs: []*url.URL{
					{
						Host: "example.com",
						Path: "/1/",
					},
					{
						Host: "example.com",
						Path: "/2/",
					},
					{
						Host: "example.com",
						Path: "/3/",
					},
				},
			},
			expectedPages: []domain.Page{
				{
					URL: url.URL{
						Host: "example.com",
						Path: "/1/",
					},
					Referrer: url.URL{
						Host: "example.com",
					},
				},
				{
					URL: url.URL{
						Host: "example.com",
						Path: "/2/",
					},
					Referrer: url.URL{
						Host: "example.com",
					},
				},
				{
					URL: url.URL{
						Host: "example.com",
						Path: "/3/",
					},
					Referrer: url.URL{
						Host: "example.com",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewController(test.givenRepo, test.givenClient, test.givenParser)

			actual, err := c.Start(context.Background(), test.givenBaseURL, test.givenBaseURL)
			if err == nil {
				t.Fatal(err)
			}

			sort.Slice(actual, func(i, j int) bool {
				return actual[i].URL.Path < actual[j].URL.Path
			})

			if !cmp.Equal(actual, test.expectedPages, cmpopts.IgnoreTypes(time.Time{})) {
				t.Fatal(cmp.Diff(actual, test.expectedPages, cmpopts.IgnoreTypes(time.Time{})))
			}
		})
	}
}

func TestController_Start_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenBaseURL  *url.URL
		givenRepo     RepositoryProvider
		givenClient   ClientProvider
		givenParser   Parser
		expectedError error
	}{
		{
			name: "given a parser fail, expect the error to be returned",
			givenBaseURL: &url.URL{
				Host: "example.com",
			},
			givenRepo: mockRepo{},
			givenClient: &mockClient{
				GivenFetchResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader(htmlBody)),
				},
			},
			givenParser: mockParser{
				GivenError: htmlparser.ErrFailedToParseHTML,
			},
			expectedError: htmlparser.ErrFailedToParseHTML,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewController(test.givenRepo, test.givenClient, test.givenParser)

			_, err := c.Start(context.Background(), test.givenBaseURL, test.givenBaseURL)
			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError))
			}
		})
	}
}

type mockRepo struct {
	GivenGetPage     domain.Page
	GivenGetError    error
	GivenInsertPage  domain.Page
	GivenInsertError error
}

func (m mockRepo) Get(_ url.URL) (domain.Page, error) {
	return m.GivenGetPage, m.GivenGetError
}

func (m mockRepo) Insert(_ domain.Page) (domain.Page, error) {
	return m.GivenInsertPage, m.GivenInsertError
}

type mockClient struct {
	GivenFetchResponse *http.Response
	GivenFetchError    error
	RequestNumber      int
}

func (m *mockClient) Fetch(_ context.Context, _ url.URL) (*http.Response, error) {
	if m.RequestNumber == 0 {
		m.RequestNumber++
		return m.GivenFetchResponse, nil
	}

	return nil, m.GivenFetchError
}

type mockParser struct {
	GivenURLs  []*url.URL
	GivenError error
}

func (m mockParser) FetchLinks(_ io.Reader, _ *url.URL) ([]*url.URL, error) {
	return m.GivenURLs, m.GivenError
}

var htmlBody = `
	<html><a href='https://example.com/1'>link1</a>
	<a href='https://example.com/2'>link2</a>
	<a href='https://example.com/3'>link3</a></html>`
