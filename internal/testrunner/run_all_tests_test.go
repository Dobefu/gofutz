package testrunner

import (
	"sync"
	"testing"
)

func TestRunAllTests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		files    map[string]File
		expected []string
	}{
		{
			name:     "no test files",
			files:    map[string]File{},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{ // nolint:exhaustruct
				files: test.files,
			}

			fullOutput := []string{}
			outputMutex := &sync.Mutex{}

			runner.RunAllTests(func(_ File) error {
				return nil
			}, func(output string) error {
				outputMutex.Lock()
				fullOutput = append(fullOutput, output)
				outputMutex.Unlock()

				return nil
			}, func() {})

			outputMutex.Lock()
			newOutput := make([]string, len(fullOutput))
			copy(newOutput, fullOutput)
			outputMutex.Unlock()

			if len(newOutput) != len(test.expected) {
				t.Fatalf(
					"expected %d output lines, got: %d",
					len(test.expected),
					len(newOutput),
				)
			}

			for i := range newOutput {
				if newOutput[i] != test.expected[i] {
					t.Fatalf(
						"expected output to be %s, got: %s",
						test.expected[i],
						newOutput[i],
					)
				}
			}
		})
	}
}

func TestSendCallbacks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		files    map[string]File
		expected map[string]bool
	}{
		{
			name: "success",
			files: map[string]File{
				"test.go": { // nolint:exhaustruct
					Name: "test.go",
				},
			},
			expected: map[string]bool{
				"test.go": true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &TestRunner{ // nolint:exhaustruct
				files: test.files,
			}

			callbackCount := 0
			callbackMutex := &sync.Mutex{}

			err := runner.sendCallbacks(
				func(_ File) error {
					callbackMutex.Lock()
					callbackCount++
					callbackMutex.Unlock()

					return nil
				},
				func() {},
				[]CoverageLine{},
				map[string]map[string]float64{},
				test.expected,
			)

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			callbackMutex.Lock()
			actualCount := callbackCount
			callbackMutex.Unlock()

			if actualCount != len(test.expected) {
				t.Fatalf(
					"expected %d callbacks, got: %d",
					len(test.expected),
					actualCount,
				)
			}
		})
	}
}
