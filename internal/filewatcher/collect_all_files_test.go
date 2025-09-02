package filewatcher

import "testing"

func TestCollectAllFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "regular test files",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			files, err := CollectAllFiles()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if len(files) == 0 {
				t.Fatalf("expected files, got none")
			}
		})
	}
}
