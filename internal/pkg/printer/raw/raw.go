package raw

import (
	"crawler/internal/domain"
	"fmt"
	"os"
)

// Printer prints and persists in ASCII for given domain.Page's.
type Printer struct {
	content []domain.Page
}

// New instantiates a Raw Printer.
func New(content []domain.Page) Printer {
	return Printer{
		content: content,
	}
}

// Print prints the given domain.Page's as ASCII.
func (c Printer) Print() (string, error) {
	return fmt.Sprint(c.content), nil
}

// Persist creates a TXT file with the given data.
func (c Printer) Persist(data string) error {
	return os.WriteFile("output.txt", []byte(data), 0600)
}
