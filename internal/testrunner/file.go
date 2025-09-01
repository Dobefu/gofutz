package testrunner

// File defines a test file.
type File struct {
	Name            string `json:"name"`
	Tests           []Test `json:"tests"`
	Code            string `json:"code"`
	HighlightedCode string `json:"highlightedCode"`
}
