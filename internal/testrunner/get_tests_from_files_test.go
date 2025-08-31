package testrunner

import (
	"fmt"
	"os"
	"testing"
)

func TestGetTestsFromFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileContents []string
		expected     []string
	}{
		{
			name: "regular test files",
			fileContents: []string{`
package testrunner
func TestGetTestsFromFiles(t *testing.T) {}
				`,
				"package testrunner",
			},
			expected: []string{"TestGetTestsFromFiles"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePaths := []string{}

			for i, fileContent := range test.fileContents {
				filePath, err := writeTestFile(t, fmt.Sprintf("%s_%d", test.name, i), fileContent)
				defer os.Remove(filePath)

				if err != nil {
					t.Fatalf("expected no error, got: %s", err.Error())
				}

				filePaths = append(filePaths, filePath)
			}

			fileTests, err := GetTestsFromFiles(filePaths)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(fileTests) != len(test.expected) {
				t.Fatalf("expected %d tests, got %d", len(test.expected), len(fileTests))
			}

			for i, fileTest := range fileTests {
				if fileTest != test.expected[i] {
					t.Fatalf("expected %s, got %s", test.expected[i], fileTest)
				}
			}
		})
	}
}

func TestGetTestsFromFilesErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileContents []string
		expected     string
	}{
		{
			name:         "empty file path",
			fileContents: []string{""},
			expected:     "file is empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePaths := []string{}

			for i, fileContent := range test.fileContents {
				filePath, err := writeTestFile(t, fmt.Sprintf("%s_%d", test.name, i), fileContent)
				defer os.Remove(filePath)

				if err != nil {
					t.Fatalf("expected no error, got: %s", err.Error())
				}

				filePaths = append(filePaths, filePath)
			}

			_, err := GetTestsFromFiles(filePaths)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != test.expected {
				t.Fatalf("expected error to be \"%s\", got \"%s\"", test.expected, err.Error())
			}
		})
	}
}
