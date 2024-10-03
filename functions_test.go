package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertCurrencyIntToString(t *testing.T) {
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
			assert.Equal(t, tc.expected, ConvertCurrencyIntToString(tc.input))
		})
	}
}

func TestConvertCurrencyStringToIntOrNil(t *testing.T) {
	tests := map[string]struct {
		input     string
		expected  int64
		expectErr bool
	}{
		"It converts 0.00 to 0":     {input: "0.00", expected: 0},
		"It converts 1.23 to 123":   {input: "1.23", expected: 123},
		"It converts -1.23 to -123": {input: "-1.23", expected: -123},
		"It errors on empty string": {input: "", expectErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ConvertCurrencyStringToInt(tc.input)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
