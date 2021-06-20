package crawler

import (
	"crawler/internal/domain"
	"crawler/storage"
	"crawler/storage/memory"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRepository_Get_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenURL     url.URL
		givenStorage StorageProvider
		expectedPage domain.Page
	}{
		{
			name:     "Given a URL, expect a page returned",
			givenURL: url.URL{Host: "https://example.com"},
			givenStorage: mockStorage{
				GivenGetPage: storage.Page{
					URL: url.URL{
						Host: "example.com",
						Path: "/test/",
					},
					Referrer:  url.URL{Host: "example.com"},
					CrawledAt: time.Date(2021, 6, 9, 11, 00, 00, 00, time.UTC),
				},
			},
			expectedPage: domain.Page{
				URL: url.URL{
					Host: "example.com",
					Path: "/test/",
				},
				Referrer:  url.URL{Host: "example.com"},
				CrawledAt: time.Date(2021, 6, 9, 11, 00, 00, 00, time.UTC),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenStorage)

			actual, err := repo.Get(test.givenURL)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedPage) {
				t.Fatal(cmp.Diff(actual, test.expectedPage))
			}
		})
	}
}

func TestRepository_Get_Fails(t *testing.T) {
	tests := []struct {
		name          string
		givenURL      url.URL
		givenStorage  StorageProvider
		expectedError error
	}{
		{
			name:     "Given a URL that has not been inserted, expect err",
			givenURL: url.URL{Host: "https://example.com"},
			givenStorage: mockStorage{
				GivenGetError: memory.ErrInvalidKey,
			},
			expectedError: memory.ErrInvalidKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenStorage)

			_, err := repo.Get(test.givenURL)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestRepository_Insert_Success(t *testing.T) {
	tests := []struct {
		name         string
		givenPage    domain.Page
		givenStorage StorageProvider
		expectedPage domain.Page
	}{
		{
			name: "Given a page, expect it inserted and returned",
			givenPage: domain.Page{
				URL: url.URL{
					Host: "example.com",
					Path: "/test",
				},
				Referrer:  url.URL{Host: "example.com"},
				CrawledAt: time.Date(2021, 6, 9, 11, 00, 00, 00, time.UTC),
			},
			givenStorage: mockStorage{},
			expectedPage: domain.Page{
				URL: url.URL{
					Host: "example.com",
					Path: "/test",
				},
				Referrer:  url.URL{Host: "example.com"},
				CrawledAt: time.Date(2021, 6, 9, 11, 00, 00, 00, time.UTC),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenStorage)

			actual, err := repo.Insert(test.givenPage)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedPage) {
				t.Fatal(cmp.Diff(actual, test.expectedPage))
			}
		})
	}
}

func TestRepository_Insert_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenPage     domain.Page
		givenStorage  StorageProvider
		expectedError error
	}{
		{
			name: "Given a page which has already been inserted, expect error to be returned",
			givenPage: domain.Page{
				URL: url.URL{
					Host: "example.com",
					Path: "/test/",
				},
				Referrer:  url.URL{Host: "example.com"},
				CrawledAt: time.Date(2021, 6, 9, 11, 00, 00, 00, time.UTC),
			},
			givenStorage: mockStorage{
				GivenInsertError: memory.ErrDuplicateKey,
			},
			expectedError: memory.ErrDuplicateKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenStorage)

			_, err := repo.Insert(test.givenPage)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockStorage struct {
	GivenGetPage     storage.Page
	GivenGetError    error
	GivenInsertError error
}

func (m mockStorage) Get(_ url.URL) (storage.Page, error) {
	return m.GivenGetPage, m.GivenGetError
}

func (m mockStorage) Insert(_ storage.Page) error {
	return m.GivenInsertError
}
