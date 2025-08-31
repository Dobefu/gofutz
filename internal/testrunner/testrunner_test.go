package testrunner

import "testing"

func TestNewTestRunner(t *testing.T) {
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
				func TestGetTestsFromFile(t *testing.T) {}
			`},
			expected: []string{"TestGetTestsFromFile"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePaths, cleanup, err := writeTestFiles(
				t,
				"TestNewTestRunner",
				test.name,
				test.fileContents,
			)

			defer cleanup()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			runner, err := NewTestRunner(filePaths)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			for i, fileTest := range runner.GetTests() {
				if fileTest == test.expected[i] {
					continue
				}

				t.Fatalf("expected %s, got %s", test.expected[i], fileTest)
			}
		})
	}
}

func TestNewTestRunnerErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileContents []string
		expected     string
	}{
		{
			name:         "invalid test files",
			fileContents: []string{""},
			expected:     "file is empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePaths, cleanup, err := writeTestFiles(
				t,
				"TestNewTestRunnerErr",
				test.name,
				test.fileContents,
			)

			defer cleanup()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			_, err = NewTestRunner(filePaths)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != test.expected {
				t.Fatalf("expected error to be \"%s\", got \"%s\"", test.expected, err.Error())
			}
		})
	}
}
