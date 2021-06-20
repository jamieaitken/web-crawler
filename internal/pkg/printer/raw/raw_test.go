package raw

import (
	"crawler/internal/domain"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestContentTypeRaw_Print_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenContent []domain.Page
		expected     string
	}{
		{
			name: "given pages, expect raw output",
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
			expected: "[{{   example.com /test/  false   } {   example.com   false   } 2021-06-10 16:00:00 +0000 UTC}]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := New(test.givenContent).Print()
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expected) {
				t.Fatal(cmp.Diff(actual, test.expected))
			}
		})
	}
}

func TestContentTypeRaw_Persist_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenContent []domain.Page
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
