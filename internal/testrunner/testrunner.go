// Package testrunner provides test runner functionality.
package testrunner

// TestRunner defines a test runner.
type TestRunner struct {
	files []string
	tests []string
}

// NewTestRunner creates a new test runner.
func NewTestRunner(files []string) (*TestRunner, error) {
	tests, err := GetTestsFromFiles(files)

	if err != nil {
		return nil, err
	}

	return &TestRunner{
		files: files,
		tests: tests,
	}, nil
}

// GetTests gets the tests.
func (t *TestRunner) GetTests() []string {
	return t.tests
}
