package printer

import (
	"crawler/internal/domain"
	"crawler/internal/pkg/printer/json"
	"crawler/internal/pkg/printer/raw"
)

// Printer instantiates the selected TypeProvider.
type Printer struct {
	typeSelected ContentType
}

// TypeProvider provides both print and persist functionality.
type TypeProvider interface {
	Print() (string, error)
	Persist(data string) error
}

// ContentType provides different types of Printers.
type ContentType string

var (
	Raw  ContentType = "raw"
	JSON ContentType = "json"
)

// New instantiates a Printer.
func New(typeSelected ContentType) Printer {
	return Printer{
		typeSelected: typeSelected,
	}
}

// Create returns the TypeProvider for the given type.
func (c Printer) Create(pages []domain.Page) TypeProvider {
	switch c.typeSelected {
	case JSON:
		return json.New(pages)
	case Raw:
		return raw.New(pages)
	default:
		return raw.New(pages)
	}
}
