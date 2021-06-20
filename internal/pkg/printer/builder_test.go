package printer

import (
	"crawler/internal/domain"
	"crawler/internal/pkg/printer/json"
	"crawler/internal/pkg/printer/raw"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

func TestPrinter_Create(t *testing.T) {
	tests := []struct {
		name                 string
		givenType            ContentType
		expectedTypeProvider TypeProvider
	}{
		{
			name:                 "given raw content type, expect raw type provider",
			givenType:            Raw,
			expectedTypeProvider: raw.Printer{},
		},
		{
			name:                 "given json content type, expect json type provider",
			givenType:            JSON,
			expectedTypeProvider: json.Printer{},
		},
		{
			name:                 "given undefined content type, default to raw type provider",
			givenType:            "test",
			expectedTypeProvider: raw.Printer{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := New(test.givenType).Create([]domain.Page{})

			if !cmp.Equal(a, test.expectedTypeProvider, cmpopts.IgnoreUnexported(raw.Printer{}, json.Printer{})) {
				t.Fatal(cmp.Diff(a, test.expectedTypeProvider, cmpopts.IgnoreUnexported(raw.Printer{}, json.Printer{})))
			}
		})
	}
}
