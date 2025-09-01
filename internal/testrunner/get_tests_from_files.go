package testrunner

// GetTestsFromFiles gets all test functions from the supplied files.
func GetTestsFromFiles(files []string) (map[string]File, error) {
	allTests := map[string]File{}

	for _, file := range files {
		tests, code, err := GetTestsFromFile(file)

		if err != nil {
			return map[string]File{}, err
		}

		allTests[file] = File{
			Name:            file,
			Tests:           tests,
			Code:            string(code),
			HighlightedCode: HighlightCode("go", string(code)),
		}
	}

	return allTests, nil
}
