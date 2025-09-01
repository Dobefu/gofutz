package testrunner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// GetTestsFromFile gets all test functions from a file.
func GetTestsFromFile(file string) ([]Test, string, error) {
	tests := []Test{}

	if file == "" {
		return []Test{}, "", fmt.Errorf("file is empty")
	}

	fileContent, err := os.ReadFile(filepath.Clean(file))

	if err != nil {
		return []Test{}, "", err
	}

	fileset := token.NewFileSet()
	node, err := parser.ParseFile(
		fileset,
		file,
		fileContent,
		parser.ParseComments,
	)

	if err != nil {
		return []Test{}, "", err
	}

	for _, declaration := range node.Decls {
		functionDeclaration, isFunctionDeclaration := declaration.(*ast.FuncDecl)

		if !isFunctionDeclaration {
			continue
		}

		if !strings.HasPrefix(functionDeclaration.Name.Name, "Test") {
			continue
		}

		tests = append(tests, Test{Name: functionDeclaration.Name.Name})
	}

	return tests, string(fileContent), nil
}
