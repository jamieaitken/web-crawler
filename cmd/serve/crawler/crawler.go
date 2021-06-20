package crawler

import (
	"context"
	"crawler/internal/crawler"
	"crawler/internal/pkg/htmlparser"
	"crawler/internal/pkg/httpclient"
	"crawler/internal/pkg/printer"
	"crawler/internal/pkg/requester"
	"crawler/internal/pkg/urlbuilder"
	"crawler/storage/memory"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

// New injects all the required dependencies for a crawler, crawls the given URL and returns the results.
func New() error {
	controller := crawler.NewController(
		crawler.NewRepository(
			memory.New(),
		),
		httpclient.New(
			&http.Client{
				Timeout: viper.GetDuration("httpTimeout"),
			},
			requester.New(),
		),
		htmlparser.New(
			urlbuilder.New(),
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	baseURL := viper.GetString("baseURL")

	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	pages, err := controller.Start(ctx, u, u)
	if err != nil {
		log.Printf("crawler errors ocuured: %v", err)
	}

	p := printer.New(printer.ContentType(viper.GetString("printerType")))

	selectedTypeContent := p.Create(pages)

	content, err := selectedTypeContent.Print()
	if err != nil {
		return err
	}

	if !viper.GetBool("persist") {
		return nil
	}

	return selectedTypeContent.Persist(content)
}
