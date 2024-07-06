package main

import (
	"testing"
)

func TestEscapeCsv(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ab c", "ab c"},
		{"a\"b", "\"a\"\"b\""},
		{"a,b", "\"a,b\""},
		{"a, b", "\"a, b\""},
		{"a,b\nc", "\"a,b\\nc\""},
		{"", ""},
		{"abc123", "abc123"},
		{"a,b,c", "\"a,b,c\""},
	}

	for _, tt := range tests {
		actual := escapeCsv(tt.input)
		if actual!= tt.expected {
			t.Errorf("escapeCsv(%q) = %q; expected %q", tt.input, actual, tt.expected)
		}
	}
}
