package htmlparser

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var (
	ErrFailedToParseHTML = errors.New("failed to parse html")
)

// HTMLParser finds all relevant hrefs for a io.Reader.
type HTMLParser struct {
	URLBuilder URLBuilder
}

// New instantiates a HTMLParser.
func New(builder URLBuilder) HTMLParser {
	return HTMLParser{
		URLBuilder: builder,
	}
}

// URLBuilder takes all found hrefs and creates url.URLs with them.
type URLBuilder interface {
	Build([]*html.Node, *url.URL) ([]*url.URL, error)
}

// FetchLinks finds all relevant hrefs for a given http.Response's body and creates url.URL for each.
func (h HTMLParser) FetchLinks(hr io.Reader, baseURL *url.URL) ([]*url.URL, error) {
	body, err := htmlquery.Parse(hr)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrFailedToParseHTML)
	}

	hrefs := htmlquery.Find(body, "//a/@href")

	return h.URLBuilder.Build(hrefs, baseURL)
}
