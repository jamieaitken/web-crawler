package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestClient_Fetch_Success(t *testing.T) {
	tests := []struct {
		name             string
		givenRequester   Requester
		givenDoer        Doer
		givenURL         url.URL
		expectedResponse *http.Response
	}{
		{
			name: "given url, expect 200 response",
			givenRequester: mockRequester{
				GivenBuildRequest: &http.Request{},
			},
			givenURL: url.URL{
				Host:   "example.com",
				Scheme: "https",
			},
			givenDoer: mockDoer{
				GivenDoResponse: &http.Response{
					StatusCode: http.StatusOK,
				},
			},
			expectedResponse: &http.Response{
				StatusCode: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := New(test.givenDoer, test.givenRequester)

			actual, err := client.Fetch(context.Background(), test.givenURL)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedResponse) {
				t.Fatal(cmp.Diff(actual, test.expectedResponse))
			}
		})
	}
}

func TestClient_Fetch_Fail(t *testing.T) {
	tests := []struct {
		name           string
		givenRequester Requester
		givenDoer      Doer
		givenURL       url.URL
		expectedError  error
	}{
		{
			name: "given failed to create request, return error",
			givenRequester: mockRequester{
				GivenBuildError: ErrFailedToBuildRequest,
			},
			givenURL: url.URL{
				Host:   "example.com",
				Scheme: "https",
			},
			givenDoer:     mockDoer{},
			expectedError: ErrFailedToBuildRequest,
		},
		{
			name:           "given failed to send request, return error",
			givenRequester: mockRequester{},
			givenURL: url.URL{
				Host:   "example.com",
				Scheme: "https",
			},
			givenDoer: mockDoer{
				GivenDoError: ErrFailedRequest,
			},
			expectedError: ErrFailedRequest,
		},
		{
			name:           "given invalid status code, return error",
			givenRequester: mockRequester{},
			givenURL: url.URL{
				Host:   "example.com",
				Scheme: "https",
			},
			givenDoer: mockDoer{
				GivenDoResponse: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			expectedError: ErrUnacceptableResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := New(test.givenDoer, test.givenRequester)

			_, err := client.Fetch(context.Background(), test.givenURL)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockRequester struct {
	GivenBuildRequest *http.Request
	GivenBuildError   error
}

func (m mockRequester) WithURL(_ url.URL) Requester {
	return m
}

func (m mockRequester) WithBody(_ io.Reader) Requester {
	return m
}

func (m mockRequester) WithMethod(_ string) Requester {
	return m
}

func (m mockRequester) Build(_ context.Context) (*http.Request, error) {
	return m.GivenBuildRequest, m.GivenBuildError
}

type mockDoer struct {
	GivenDoResponse *http.Response
	GivenDoError    error
}

func (m mockDoer) Do(_ *http.Request) (*http.Response, error) {
	return m.GivenDoResponse, m.GivenDoError
}
