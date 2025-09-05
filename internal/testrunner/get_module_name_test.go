package testrunner

import "testing"

func TestGetModuleNameErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "cannot find go.mod",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			moduleName := GetModuleName()

			if moduleName != test.expected {
				t.Fatalf(
					"expected module name to be \"%s\", got \"%s\"",
					test.expected,
					moduleName,
				)
			}
		})
	}
}
