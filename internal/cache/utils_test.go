package cache

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Bytes", 500, "500 B"},
		{"Kilobytes", 1024, "1.0 KB"},
		{"Kilobytes fractional", 1536, "1.5 KB"},
		{"Megabytes", 1024 * 1024, "1.0 MB"},
		{"Gigabytes", 1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBytes(tt.input); got != tt.expected {
				t.Errorf("FormatBytes(%d) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatCount(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Small number", 123, "123"},
		{"Thousand", 1000, "1,000"},
		{"Ten thousand", 10000, "10,000"},
		{"Million", 1000000, "1,000,000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatCount(tt.input); got != tt.expected {
				t.Errorf("FormatCount(%d) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
