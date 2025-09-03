package testrunner

import (
	"fmt"
	"strings"
)

// GetFunctionsFromFiles gets all functions from the supplied files.
func GetFunctionsFromFiles(files []string) (map[string]File, error) {
	moduleName := GetModuleName()

	allFunctions := map[string]File{}
	testFiles := map[string][]Function{}

	for _, file := range files {
		functions, code, err := GetFunctionsFromFile(file)

		if err != nil {
			return map[string]File{}, err
		}

		if strings.HasSuffix(file, "_test.go") {
			sourceFile := strings.Replace(file, "_test.go", ".go", 1)
			testFiles[sourceFile] = functions

			continue
		}

		if moduleName != "" {
			file = fmt.Sprintf("%s/%s", moduleName, file)
		}

		allFunctions[file] = File{
			Name:            file,
			Functions:       functions,
			Code:            string(code),
			HighlightedCode: HighlightCode("go", string(code)),
			Coverage:        -1,
			CoveredLines:    []Line{},
		}
	}

	for sourceFile, functions := range testFiles {
		if moduleName != "" {
			sourceFile = fmt.Sprintf("%s/%s", moduleName, sourceFile)
		}

		file, hasFile := allFunctions[sourceFile]

		if !hasFile {
			continue
		}

		file.Functions = append(file.Functions, functions...)
		allFunctions[sourceFile] = file
	}

	return allFunctions, nil
}
