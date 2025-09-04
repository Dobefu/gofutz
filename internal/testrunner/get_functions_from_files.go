package testrunner

import (
	"fmt"
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

		if moduleName != "" {
			file = fmt.Sprintf("%s/%s", moduleName, file)
		}

		var status TestStatus = TestStatusPending

		if len(functions) == 0 {
			status = TestStatusNoCodeToCover
		}

		allFunctions[file] = File{
			Name:            file,
			Functions:       functions,
			Code:            string(code),
			HighlightedCode: HighlightCode("go", string(code)),
			Status:          status,
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
