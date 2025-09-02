package testrunner

import "log/slog"

// RunAllTests runs all tests.
func (t *TestRunner) RunAllTests(testCompleteCallback func(test Test) error) {
	for _, file := range t.tests {
		for _, test := range file.Tests {
			go func() {
				test.Run()

				err := testCompleteCallback(test)

				if err != nil {
					slog.Error(err.Error())
				}
			}()
		}
	}
}
