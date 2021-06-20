package memory

import (
	"crawler/storage"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMemory_Get_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenURL     url.URL
		givenMemory  *Memory
		expectedPage storage.Page
	}{
		{
			name: "given a URL that exists, return page",
			givenURL: url.URL{
				Host: "example.com",
			},
			givenMemory: &Memory{
				pages: map[url.URL]storage.Page{
					{Host: "example.com"}: {
						Referrer: url.URL{
							Host: "example.com",
						},
					},
				},
			},
			expectedPage: storage.Page{
				Referrer: url.URL{
					Host: "example.com",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.givenMemory.Get(test.givenURL)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedPage) {
				t.Fatal(cmp.Diff(actual, test.expectedPage))
			}
		})
	}
}

func TestMemory_Get_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenURL      url.URL
		givenMemory   *Memory
		expectedError error
	}{
		{
			name: "given a URL that does not exist, return error",
			givenURL: url.URL{
				Host: "example.com",
			},
			givenMemory:   &Memory{},
			expectedError: ErrInvalidKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.givenMemory.Get(test.givenURL)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestMemory_Insert_Success(t *testing.T) {
	tests := []struct {
		name        string
		givenPage   storage.Page
		givenMemory *Memory
	}{
		{
			name: "given a page that does not exist, insert",
			givenPage: storage.Page{
				URL: url.URL{
					Host: "example.com",
					Path: "/test/",
				},
				Referrer: url.URL{
					Host: "example.com",
				},
				CrawledAt: time.Time{},
			},
			givenMemory: &Memory{
				pages: map[url.URL]storage.Page{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.givenMemory.Insert(test.givenPage)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMemory_Insert_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenPage     storage.Page
		givenMemory   *Memory
		expectedError error
	}{
		{
			name: "given a page that exists, return error",
			givenPage: storage.Page{
				URL: url.URL{
					Host: "example.com",
					Path: "/test/",
				},
				Referrer: url.URL{
					Host: "example.com",
				},
				CrawledAt: time.Time{},
			},
			givenMemory: &Memory{
				pages: map[url.URL]storage.Page{
					{Host: "example.com", Path: "/test/"}: {
						URL: url.URL{
							Host: "example.com",
							Path: "/test/",
						},
						Referrer: url.URL{
							Host: "example.com",
						},
						CrawledAt: time.Time{},
					},
				},
			},
			expectedError: ErrDuplicateKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.givenMemory.Insert(test.givenPage)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}
