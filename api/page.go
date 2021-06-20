package api

import (
	"time"
)

// Page shows which URLs were found on a given URL.
type Page struct {
	URL       URL       `json:"url"`
	Referrer  URL       `json:"referrer"`
	CrawledAt time.Time `json:"crawledAt"`
}

// URL is a JSON representation of url.URL.
type URL struct {
	Scheme      string `json:"scheme"`
	Opaque      string `json:"opaque"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	RawPath     string `json:"rawPath"`
	ForceQuery  bool   `json:"forceQuery"`
	RawQuery    string `json:"rawQuery"`
	Fragment    string `json:"fragment"`
	RawFragment string `json:"rawFragment"`
}
