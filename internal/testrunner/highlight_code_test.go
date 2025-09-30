package testrunner

import (
	"regexp"
	"testing"
)

func TestHighlightCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		language string
		code     string
		expected string
	}{
		{
			name:     "go code",
			language: "go",
			code:     "package main",
			expected: "<.+?>package</.+?><.+?>main</.+?>",
		},
		{
			name:     "empty language",
			language: "",
			code:     "package main",
			expected: "package main",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := HighlightCode(test.language, test.code)
			isMatch, err := regexp.MatchString(test.expected, result)

			if err != nil {
				t.Fatalf("expected no error, got \"%s\"", err)
			}

			if !isMatch {
				t.Fatalf(
					"expected output to match \"%s\", got \"%s\"",
					test.expected,
					result,
				)
			}
		})
	}
}
