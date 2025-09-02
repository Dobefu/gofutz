package testrunner

import "testing"

func TestNewTestRunner(t *testing.T) {
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
				func TestGetTestsFromFile(t *testing.T) {}
			`},
			expected: []Test{
				{
					Name: "TestGetTestsFromFile",
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

			for _, fileTest := range runner.GetFiles() {
				if fileTest.Tests[0].Name == test.expected[0].Name {
					return
				}
			}

			t.Fatalf(
				"expected \"%s\", got \"%s\"",
				test.expected[0].Name,
				runner.GetFiles()[filePaths[0]].Tests[0].Name,
			)
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
