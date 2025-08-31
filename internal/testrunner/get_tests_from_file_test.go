package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(
	t *testing.T,
	name string,
	fileContent string,
) (string, func(), error) {
	t.Helper()

	filePath := ""

	if fileContent != "" {
		filePath = filepath.Join(
			os.TempDir(),
			fmt.Sprintf("test_%s.go", name),
		)

		if fileContent != "-" {
			err := os.WriteFile(filePath, []byte(fileContent), 0600)

			if err != nil {
				return "", func() {}, err
			}
		}
	}

	return filePath, func() { _ = os.Remove(filePath) }, nil
}

func TestGetTestsFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		fileContent string
		expected    []string
	}{
		{
			name: "regular test file",
			fileContent: `
package testrunner
func TestGetTestsFromFile(t *testing.T) {
				`,
			expected: []string{"TestGetTestsFromFile"},
		},
		{
			name: "generic test function",
			fileContent: `
package testrunner
func TestGetTestsFromFile[T any](t *testing.T) {}
				`,
			expected: []string{"TestGetTestsFromFile"},
		},
		{
			name: "invalid test function syntax",
			fileContent: `
package testrunner
func TestGetTestsFromFile
				`,
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePath, cleanup, err := writeTestFile(t, test.name, test.fileContent)
			defer cleanup()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			fileTests, err := GetTestsFromFile(filePath)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			for i, fileTest := range fileTests {
				if fileTest != test.expected[i] {
					t.Fatalf("expected %s, got %s", test.expected[i], fileTest)
				}
			}
		})
	}
}

func TestGetTestsFromFileErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		fileContent string
		expected    string
	}{
		{
			name:        "empty file path",
			fileContent: "",
			expected:    "file is empty",
		},
		{
			name:        "no file found",
			fileContent: "-",
			expected:    "no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePath, cleanup, err := writeTestFile(t, test.name, test.fileContent)
			defer cleanup()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			_, err = GetTestsFromFile(filePath)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !strings.Contains(err.Error(), test.expected) {
				t.Errorf(
					"expected error to contain \"%s\", got \"%s\"",
					test.expected,
					err.Error(),
				)
			}
		})
	}
}
