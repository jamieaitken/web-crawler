package httpclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Errors returned by the Client.
var (
	ErrFailedRequest        = errors.New("failed to get response")
	ErrFailedToBuildRequest = errors.New("failed to build request")
	ErrUnacceptableResponse = errors.New("invalid status code received")
)

// Doer sends a http.Request and returns a http.Response.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Requester provides a way of creating http.Request's.
type Requester interface {
	WithURL(url url.URL) Requester
	WithBody(body io.Reader) Requester
	WithMethod(method string) Requester
	Build(ctx context.Context) (*http.Request, error)
}

// Client builds a request for a given url.URL and returns the response.
type Client struct {
	Doer      Doer
	Requester Requester
}

// New instantiates a Client.
func New(doer Doer, requester Requester) *Client {
	return &Client{
		Doer:      doer,
		Requester: requester,
	}
}

// Fetch performs a http.MethodGet for a given url.URL.
func (c *Client) Fetch(ctx context.Context, u url.URL) (*http.Response, error) {
	req, err := c.Requester.
		WithMethod(http.MethodGet).
		WithURL(u).
		Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrFailedToBuildRequest)
	}

	res, err := c.Doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrFailedRequest)
	}

	if !isAcceptableStatus(res.StatusCode) {
		return nil, ErrUnacceptableResponse
	}

	return res, nil
}

func isAcceptableStatus(status int) bool {
	return status == http.StatusOK
}
