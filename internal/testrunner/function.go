package testrunner

// TestStatus defines the status of a test.
type TestStatus int

const (
	// TestStatusPending represents a test that is pending.
	TestStatusPending = iota
	// TestStatusRunning represents a test that is running.
	TestStatusRunning
	// TestStatusPassed represents a test that has passed.
	TestStatusPassed
	// TestStatusFailed represents a test that has failed.
	TestStatusFailed
	// TestStatusNoTests represents a file with no test functions.
	TestStatusNoTests
)

// Function defines a single function.
type Function struct {
	Name   string     `json:"name"`
	Result TestResult `json:"result"`
}

// TestResult defines the result of a test.
type TestResult struct {
	Coverage float64 `json:"coverage"`
}

// Line defines a line of code in a coverage report.
type Line struct {
	Number             int `json:"number"`
	StartLine          int `json:"startLine"`
	StartColumn        int `json:"startColumn"`
	EndLine            int `json:"endLine"`
	EndColumn          int `json:"endColumn"`
	ExecutionCount     int `json:"executionCount"`
	NumberOfStatements int `json:"numberOfStatements"`
}
