package testrunner

// File defines a test file.
type File struct {
	Name            string     `json:"name"`
	Functions       []Function `json:"functions"`
	Code            string     `json:"code"`
	HighlightedCode string     `json:"highlightedCode"`
	Coverage        float64    `json:"coverage"`
	CoveredLines    []Line     `json:"coveredLines"`
}
