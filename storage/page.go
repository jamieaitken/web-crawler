package storage

import (
	"net/url"
	"time"
)

// Page is the storage representation of domain.Page.
type Page struct {
	URL       url.URL
	Referrer  url.URL
	CrawledAt time.Time
}
