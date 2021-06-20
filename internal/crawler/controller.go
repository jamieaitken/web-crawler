package crawler

import (
	"context"
	"crawler/internal/domain"
	"crawler/storage/memory"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// Controller fetches all available URLs and returns them.
type Controller struct {
	Repository RepositoryProvider
	Client     ClientProvider
	Parser     Parser
	err        chan error
	timeout    *time.Timer
	page       chan domain.Page
	foundURLS  chan *url.URL
}

// NewController instantiates a Controller.
func NewController(repo RepositoryProvider, client ClientProvider, parser Parser) Controller {
	return Controller{
		Repository: repo,
		Client:     client,
		Parser:     parser,
		err:        make(chan error),
		page:       make(chan domain.Page),
		foundURLS:  make(chan *url.URL),
		timeout:    time.NewTimer(time.Second * 10),
	}
}

// RepositoryProvider gives access to a storage interface.
type RepositoryProvider interface {
	Get(url url.URL) (domain.Page, error)
	Insert(page domain.Page) (domain.Page, error)
}

// ClientProvider gives the ability to perform HTTP Requests.
type ClientProvider interface {
	Fetch(ctx context.Context, url url.URL) (*http.Response, error)
}

// Parser will find all the desired URLs for a given http.Response.
type Parser interface {
	FetchLinks(html io.Reader, baseURL *url.URL) ([]*url.URL, error)
}

// Start initiates the crawler and returns the Pages found and any errors that happened.
func (c *Controller) Start(ctx context.Context, baseURL, targetURL *url.URL) ([]domain.Page, error) {
	var pageResults []domain.Page
	var errs error

	c.crawl(ctx, baseURL, targetURL)

	for {
		select {
		case page := <-c.page:
			pageResults = append(pageResults, page)
			log.Infof("received page from channel: %v", page.URL)
			c.timeout.Reset(time.Second * 10)
		case foundURL := <-c.foundURLS:
			log.Infof("received URL from channel: %v", foundURL.String())
			go c.crawl(ctx, baseURL, foundURL)
		case err := <-c.err:
			errs = fmt.Errorf("%v: %w", errs, err)
			log.Infof("received err from channel: %v", err)
		case <-c.timeout.C:
			return pageResults, errs
		}
	}
}

func (c *Controller) crawl(ctx context.Context, baseURL, targetURL *url.URL) {
	res, err := c.Client.Fetch(ctx, *targetURL)
	if err != nil {
		log.Errorf("fetch error for %v", targetURL)

		go func() {
			c.err <- fmt.Errorf("%v: %w", targetURL, err)
		}()

		return
	}

	links, err := c.Parser.FetchLinks(res.Body, baseURL)
	if err != nil {
		log.Errorf("create links for %v", targetURL)

		go func() {
			c.err <- err
		}()

		return
	}

	for _, link := range links {
		_, err = c.Repository.Get(*link)
		if !errors.Is(err, memory.ErrInvalidKey) {
			continue
		}

		page := domain.Page{
			Referrer:  *targetURL,
			URL:       *link,
			CrawledAt: time.Now().UTC(),
		}

		_, err := c.Repository.Insert(page)
		if err != nil {
			log.Infof("repo for %v", targetURL)
			continue
		}

		link := link
		go func() {
			c.foundURLS <- link
			log.Infof("a url has been sent to the channel: %v", link.String())

			c.page <- page
			log.Infof("a page has been sent to the channel: %v", page.URL)
		}()
	}

	log.Infof("all URLs have been crawled for %v", targetURL)
}
