package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetFunctionsFromFiles gets all functions from the supplied files.
func GetFunctionsFromFiles(files []string) (map[string]File, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return map[string]File{}, err
	}

	moduleName := GetModuleName(filepath.Join(cwd, "go.mod"))

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

		if len(functions) == 0 {
			continue
		}

		allFunctions[file] = File{
			Name:            file,
			Functions:       functions,
			Code:            code,
			HighlightedCode: HighlightCode("go", code),
			Status:          TestStatusPending,
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
