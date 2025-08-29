package templates

import (
	"testing"
)

func TestGetPageTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		page     string
		expected string
	}{
		{
			name:     "index page template",
			page:     "index",
			expected: "index.gohtml",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tmpl, err := GetPageTemplate(test.page)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if tmpl == nil {
				t.Fatalf("expected template, got nil")
			}

			if tmpl.Name() != test.expected {
				t.Errorf(
					"expected template name to be \"%s\", got \"%s\"",
					test.expected,
					tmpl.Name(),
				)
			}
		})
	}
}

func TestGetPageTemplateErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		page     string
		expected string
	}{
		{
			name:     "nonexistent page template",
			page:     "bogus",
			expected: "template: pattern matches no files: `pages/bogus.gohtml`",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := GetPageTemplate(test.page)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, err.Error())
			}
		})
	}
}
