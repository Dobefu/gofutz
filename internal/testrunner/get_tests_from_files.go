package testrunner

import (
	"fmt"
	"strings"
)

// GetTestsFromFiles gets all test functions from the supplied files.
func GetTestsFromFiles(files []string) (map[string]File, error) {
	moduleName := GetModuleName()

	allTests := map[string]File{}
	testFiles := map[string][]Test{}

	for _, file := range files {
		tests, code, err := GetTestsFromFile(file)

		if err != nil {
			return map[string]File{}, err
		}

		if strings.HasSuffix(file, "_test.go") {
			sourceFile := strings.Replace(file, "_test.go", ".go", 1)
			testFiles[sourceFile] = tests

			continue
		}

		if moduleName != "" {
			file = fmt.Sprintf("%s/%s", moduleName, file)
		}

		allTests[file] = File{
			Name:            file,
			Tests:           tests,
			Code:            string(code),
			HighlightedCode: HighlightCode("go", string(code)),
			Coverage:        0,
		}
	}

	for sourceFile, tests := range testFiles {
		if moduleName != "" {
			sourceFile = fmt.Sprintf("%s/%s", moduleName, sourceFile)
		}

		file, hasFile := allTests[sourceFile]

		if !hasFile {
			continue
		}

		file.Tests = append(file.Tests, tests...)
		allTests[sourceFile] = file
	}

	return allTests, nil
}
