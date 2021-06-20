package crawler

import (
	"crawler/internal/domain"
	"crawler/storage"
	"net/url"
)

// StorageProvider provides access to a storage medium.
type StorageProvider interface {
	Get(url url.URL) (storage.Page, error)
	Insert(page storage.Page) error
}

// Repository adapts between domain and storage models.
type Repository struct {
	Storage StorageProvider
}

// NewRepository instantiates a Repository.
func NewRepository(s StorageProvider) Repository {
	return Repository{
		Storage: s,
	}
}

// Get a domain.Page for a given url.URL.
func (r Repository) Get(u url.URL) (domain.Page, error) {
	page, err := r.Storage.Get(u)
	if err != nil {
		return domain.Page{}, err
	}

	return adaptStorageToDomain(page), nil
}

// Insert a domain.Page into storage.
func (r Repository) Insert(page domain.Page) (domain.Page, error) {
	err := r.Storage.Insert(adaptStorageFromDomain(page))
	if err != nil {
		return domain.Page{}, err
	}

	return page, nil
}

func adaptStorageToDomain(page storage.Page) domain.Page {
	return domain.Page{
		URL:       page.URL,
		Referrer:  page.Referrer,
		CrawledAt: page.CrawledAt,
	}
}

func adaptStorageFromDomain(page domain.Page) storage.Page {
	return storage.Page{
		URL:       page.URL,
		Referrer:  page.Referrer,
		CrawledAt: page.CrawledAt,
	}
}
