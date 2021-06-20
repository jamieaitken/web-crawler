package json

import (
	"crawler/api"
	"crawler/internal/domain"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestContentTypeJSON_Print_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenContent []domain.Page
		expected     []api.Page
	}{
		{
			name: "given pages, expect JSON output",
			givenContent: []domain.Page{
				{
					URL: url.URL{
						Host: "example.com",
						Path: "/test/",
					},
					Referrer: url.URL{
						Host: "example.com",
					},
					CrawledAt: time.Date(2021, 06, 10, 16, 00, 00, 00, time.UTC),
				},
			},
			expected: []api.Page{
				{
					URL: api.URL{
						Host: "example.com",
						Path: "/test/",
					},
					Referrer: api.URL{
						Host: "example.com",
					},
					CrawledAt: time.Date(2021, 06, 10, 16, 00, 00, 00, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := New(test.givenContent).Print()
			if err != nil {
				t.Fatal(err)
			}

			var a []api.Page

			err = json.Unmarshal([]byte(actual), &a)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(a, test.expected) {
				t.Fatal(cmp.Diff(a, test.expected))
			}
		})
	}
}

func TestContentTypeJSON_Persist_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenContent []domain.Page
		expected     string
	}{
		{
			name: "given pages, expect to be persisted",
			givenContent: []domain.Page{
				{
					Referrer: url.URL{
						Host: "example.com",
					},
					CrawledAt: time.Date(2021, 06, 10, 16, 00, 00, 00, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			printer := New(test.givenContent)

			content, err := printer.Print()
			if err != nil {
				t.Fatal(err)
			}

			err = printer.Persist(content)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
