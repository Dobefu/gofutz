package testrunner

import (
	"fmt"
	"testing"
)

func writeTestFiles(
	t *testing.T,
	prefix string,
	name string,
	fileContents []string,
) ([]string, func(), error) {
	t.Helper()

	filePaths := []string{}
	cleanups := []func(){}

	for i, fileContent := range fileContents {
		filePath, cleanup, err := writeTestFile(
			t,
			prefix,
			fmt.Sprintf("%s_%d", name, i),
			fileContent,
		)

		if err != nil {
			return []string{}, func() {}, err
		}

		filePaths = append(filePaths, filePath)
		cleanups = append(cleanups, cleanup)
	}

	return filePaths, func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}, nil
}

func TestGetTestsFromFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileContents []string
		expected     []Test
	}{
		{
			name: "regular test files",
			fileContents: []string{`
package testrunner
func TestGetTestsFromFiles(t *testing.T) {}
				`,
				"package testrunner",
			},
			expected: []Test{
				{
					Name: "TestGetTestsFromFiles",
					Result: TestResult{
						Status:       TestStatusPending,
						Output:       []string{},
						Coverage:     0,
						CoveredLines: []Line{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePaths, cleanup, err := writeTestFiles(
				t,
				"TestGetTestsFromFiles",
				test.name,
				test.fileContents,
			)

			defer cleanup()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			fileTests, err := GetTestsFromFiles(filePaths)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			for i, expected := range test.expected {
				if fileTests[filePaths[i]].Tests[i].Name != expected.Name {
					t.Fatalf(
						"expected %s, got %s",
						expected.Name,
						fileTests[filePaths[i]].Tests[i].Name,
					)
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

			filePaths, cleanup, err := writeTestFiles(
				t,
				"TestGetTestsFromFiles",
				test.name,
				test.fileContents,
			)

			defer cleanup()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			_, err = GetTestsFromFiles(filePaths)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != test.expected {
				t.Fatalf(
					"expected error to be \"%s\", got \"%s\"",
					test.expected,
					err.Error(),
				)
			}
		})
	}
}
