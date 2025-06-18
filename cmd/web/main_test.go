package main

import "testing"

// *_test.go is a test file for *
// For testing on terminal, run 'go test' or 'go test -v'
// Run 'go test -coverage' to get coverage of test
// Run 'go test -coverprofile=coverage.out && go tool cover -html=coverage.out' runs Go tests, collects code coverage data, and then displays that coverage in your browser as an HTML report.

// Testing main function
func TestRun(t *testing.T) {
	_, err := run()
	if err != nil {
		t.Error("failed run()")
	}
}
