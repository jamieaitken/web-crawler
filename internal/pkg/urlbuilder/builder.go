package urlbuilder

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// ErrFailedToParseURL is returned if it's an invalid url.URL.
var (
	ErrFailedToParseURL = errors.New("failed to parse url")
)

// Builder takes the found hrefs from a http.Response body and builds desired url.URL's.
type Builder struct {
}

// New instantiates a Builder.
func New() Builder {
	return Builder{}
}

// Build takes the found hrefs from a http.Response body and builds desired url.URL's.
func (b Builder) Build(nodes []*html.Node, baseURL *url.URL) ([]*url.URL, error) {
	urls := make(map[url.URL]*url.URL)

	for _, node := range nodes {
		n := htmlquery.SelectAttr(node, "href")

		u, err := url.Parse(n)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, ErrFailedToParseURL)
		}

		if u.Hostname() != baseURL.Hostname() && u.IsAbs() {
			continue
		}

		// This is needed as we don't want to crawl pages that are protected by Cloudflare.
		if strings.Contains(u.Path, "cdn-cgi") {
			continue
		}

		u = cleanUpURL(u, baseURL)

		_, ok := urls[*u]
		if !ok {
			urls[*u] = u
		}
	}

	return mapToArray(urls), nil
}

func cleanUpURL(u, baseURL *url.URL) *url.URL {
	u.RawQuery = ""

	u = baseURL.ResolveReference(u)

	if !strings.HasSuffix(u.Path, "/") {
		u.Path = fmt.Sprintf("%v/", u.Path)
	}

	return u
}

func mapToArray(urls map[url.URL]*url.URL) []*url.URL {
	var urlArray []*url.URL

	for _, u := range urls {
		urlArray = append(urlArray, u)
	}

	return urlArray
}
