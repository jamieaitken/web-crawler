package json

import (
	"crawler/api"
	"crawler/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
)

// ErrFailedToMarshal is returned if the marshaller fails.
var (
	ErrFailedToMarshal = errors.New("failed to marshal content")
)

// Printer prints and persists in JSON for given domain.Page's.
type Printer struct {
	content []domain.Page
}

// New instantiates a JSON Printer.
func New(content []domain.Page) Printer {
	return Printer{
		content: content,
	}
}

// Print adapts the given domain.Page's to api.Page's.
func (c Printer) Print() (string, error) {
	presentationPages := make([]api.Page, len(c.content))

	for i := range c.content {
		presentationPages[i] = api.Page{
			URL:       adaptPresentationURLFromDomain(c.content[i].URL),
			Referrer:  adaptPresentationURLFromDomain(c.content[i].Referrer),
			CrawledAt: c.content[i].CrawledAt,
		}
	}

	b, err := json.MarshalIndent(presentationPages, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%v: %w", err, ErrFailedToMarshal)
	}

	return string(b), nil
}

// Persist creates a JSON file with the given data.
func (c Printer) Persist(data string) error {
	return os.WriteFile("output.json", []byte(data), 0600)
}

func adaptPresentationURLFromDomain(u url.URL) api.URL {
	return api.URL{
		Scheme:      u.Scheme,
		Opaque:      u.Opaque,
		Host:        u.Host,
		Path:        u.Path,
		RawPath:     u.RawPath,
		ForceQuery:  u.ForceQuery,
		RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
	}
}
