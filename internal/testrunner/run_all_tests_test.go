package testrunner

import (
	"errors"
	"os"
	"sync"
	"testing"
	"time"
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

func TestRunAllTestsCancel(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		files: map[string]File{},
	}

	var isCompletionCalled bool
	var mu sync.Mutex

	runner.RunAllTests(
		func(_ File) error { return nil },
		func(_ string) error { return nil },
		func() {
			mu.Lock()
			isCompletionCalled = true
			mu.Unlock()
		},
	)

	runner.StopTests()
	time.Sleep(100 * time.Millisecond)

	if !isCompletionCalled {
		t.Fatalf("expected completion callback to be called")
	}
}

func TestHandleTestFailure(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		files:      map[string]File{},
		cancelFunc: func() {},
	}

	runner.handleTestFailure(
		[]byte{},
		errors.New("test error"),
		func(_ File) error { return nil },
		func(_ string) error { return nil },
		func() {},
	)

	if runner.cancelFunc != nil {
		t.Fatalf("expected cancel func to be nil, got: %T", runner.cancelFunc)
	}
}

func TestHandleTestSuccess(t *testing.T) {
	t.Parallel()

	coverageFile, err := os.CreateTemp("", "coverage_TestHandleTestSuccess.out")

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	defer func() { _ = os.Remove(coverageFile.Name()) }()

	var isCompletionCalled bool
	var mu sync.Mutex

	runner := &TestRunner{ // nolint:exhaustruct
		files:      map[string]File{},
		cancelFunc: func() {},
	}

	runner.handleTestSuccess(
		[]byte{},
		coverageFile.Name(),
		time.Now(),
		func(_ File) error { return nil },
		func(_ string) error { return nil },
		func() {
			mu.Lock()
			isCompletionCalled = true
			mu.Unlock()
		},
	)

	if runner.cancelFunc != nil {
		t.Fatalf("expected cancel func to be nil, got: %T", runner.cancelFunc)
	}

	if !isCompletionCalled {
		t.Fatalf("expected completion callback to be called")
	}
}

func TestHandleTestSuccessErr(t *testing.T) {
	t.Parallel()

	var isCompletionCalled bool
	var mu sync.Mutex

	runner := &TestRunner{ // nolint:exhaustruct
		files:      map[string]File{},
		cancelFunc: func() {},
	}

	runner.handleTestSuccess(
		[]byte{},
		"/nonexistent/coverage.out",
		time.Now(),
		func(_ File) error { return nil },
		func(_ string) error { return nil },
		func() {
			mu.Lock()
			isCompletionCalled = true
			mu.Unlock()
		},
	)

	if isCompletionCalled {
		t.Fatalf("expected completion callback to not be called")
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
