package domain

import (
	"net/url"
	"time"
)

// Page is the domain representation of a crawled web-page.
type Page struct {
	URL       url.URL
	Referrer  url.URL
	CrawledAt time.Time
}
