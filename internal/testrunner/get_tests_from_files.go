package testrunner

// GetTestsFromFiles gets all test functions from the supplied files.
func GetTestsFromFiles(files []string) ([]string, error) {
	allTests := []string{}

	for _, file := range files {
		tests, err := GetTestsFromFile(file)

		if err != nil {
			return []string{}, err
		}

		allTests = append(allTests, tests...)
	}

	return allTests, nil
}
