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
)

// Test defines a single test.
type Test struct {
	Name   string     `json:"name"`
	Result TestResult `json:"result"`
}

// TestResult defines the result of a test.
type TestResult struct {
	Status       TestStatus `json:"status"`
	Output       []string   `json:"output"`
	Coverage     float64    `json:"coverage"`
	CoveredLines []Line     `json:"coveredLines"`
}

// Line defines a line of code in a coverage report.
type Line struct {
	Number         int `json:"number"`
	ExecutionCount int `json:"executionCount"`
}
