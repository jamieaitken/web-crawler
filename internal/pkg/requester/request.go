package requester

import (
	"context"
	"crawler/internal/pkg/httpclient"
	"io"
	"net/http"
	"net/url"
)

// Requester builds a http.Request for the given inputs.
type Requester struct {
	url    url.URL
	body   io.Reader
	method string
}

// New instantiates a Requester.
func New() *Requester {
	return &Requester{}
}

// WithURL sets the url of the http.Request.
func (r *Requester) WithURL(u url.URL) httpclient.Requester {
	r.url = u

	return r
}

// WithBody sets the body payload of the http.Request.
func (r *Requester) WithBody(body io.Reader) httpclient.Requester {
	r.body = body

	return r
}

// WithMethod sets the type of http.Request.
func (r *Requester) WithMethod(method string) httpclient.Requester {
	r.method = method

	return r
}

// Build takes the given inputs and creates http.Request with the given context.Context.
func (r *Requester) Build(ctx context.Context) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, r.method, r.url.String(), r.body)
}
