package serve

import (
	"crawler/cmd/serve/crawler"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// NewCmd associated the serve command with the instantiation of a Crawler.
func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "serve instantiates a crawler",
		Run:   Run,
	}
}

// Run instantiates a Crawler.
func Run(_ *cobra.Command, _ []string) {
	err := crawler.New()
	if err != nil {
		log.Printf("crawler err: %v", err)
	}
}
