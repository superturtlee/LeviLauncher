package main

import (
	"testing"
)

func TestEnsureMsixvcFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "download.msixvc",
		},
		{
			name:     "filename without extension",
			input:    "test",
			expected: "test.msixvc",
		},
		{
			name:     "filename with extension",
			input:    "test.msixvc",
			expected: "test.msixvc",
		},
		{
			name:     "filename with uppercase extension",
			input:    "test.MSIXVC",
			expected: "test.MSIXVC",
		},
		{
			name:     "filename with path separators",
			input:    "path/to/file",
			expected: "path_to_file.msixvc",
		},
		{
			name:     "filename with backslashes",
			input:    "path\\to\\file",
			expected: "path_to_file.msixvc",
		},
		{
			name:     "filename with spaces",
			input:    "  test.msixvc  ",
			expected: "test.msixvc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureMsixvcFilename(tt.input)
			if result != tt.expected {
				t.Errorf("ensureMsixvcFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDeriveFilename(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "URL with filename parameter",
			url:      "http://example.com/download?filename=test",
			expected: "test.msixvc",
		},
		{
			name:     "URL with filename in path",
			url:      "http://example.com/downloads/test.msixvc",
			expected: "test.msixvc",
		},
		{
			name:     "URL with filename without extension",
			url:      "http://example.com/downloads/1.21.50.07",
			expected: "1.21.50.07.msixvc",
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "download.msixvc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveFilename(tt.url)
			if result != tt.expected {
				t.Errorf("deriveFilename(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}
