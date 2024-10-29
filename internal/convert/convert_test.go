package convert

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrencyIntToString(t *testing.T) {
	tests := map[string]struct {
		input    int64
		expected string
	}{
		"It converts 0 to 0.00":     {input: 0, expected: "0.00"},
		"It converts 100 to 1.23":   {input: 123, expected: "1.23"},
		"It converts -100 to -1.23": {input: -123, expected: "-1.23"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, CurrencyIntToString(tc.input))
		})
	}
}

func TestCurrencyStringToIntOrNil(t *testing.T) {
	tests := map[string]struct {
		input     string
		expected  sql.NullInt64
		expectErr bool
	}{
		"It converts 0.00 to 0":              {input: "0.00", expected: sql.NullInt64{Int64: 0, Valid: true}},
		"It converts 1.23 to 123":            {input: "1.23", expected: sql.NullInt64{Int64: 123, Valid: true}},
		"It converts -1.23 to -123":          {input: "-1.23", expected: sql.NullInt64{Int64: -123, Valid: true}},
		"It returns empty string as invalid": {input: "", expected: sql.NullInt64{Valid: false}},
    "On bad data it invalidates": {input: "bad data", expected: sql.NullInt64{Valid: false}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := CurrencyStringToNullInt64(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
