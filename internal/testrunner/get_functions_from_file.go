package testrunner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// GetFunctionsFromFile gets all functions from a file.
func GetFunctionsFromFile(file string) ([]Function, string, error) {
	functions := []Function{}

	if file == "" {
		return []Function{}, "", fmt.Errorf("file is empty")
	}

	fileContent, err := os.ReadFile(filepath.Clean(file))

	if err != nil {
		return []Function{}, "", err
	}

	fileset := token.NewFileSet()
	node, err := parser.ParseFile(
		fileset,
		file,
		fileContent,
		parser.ParseComments,
	)

	if err != nil {
		return []Function{}, "", err
	}

	for _, declaration := range node.Decls {
		functionDeclaration, isFunctionDeclaration := declaration.(*ast.FuncDecl)

		if !isFunctionDeclaration {
			continue
		}

		functions = append(functions, Function{
			Name: functionDeclaration.Name.Name,
			Result: TestResult{
				Status:   TestStatusPending,
				Coverage: -1,
			},
		})
	}

	return functions, string(fileContent), nil
}
