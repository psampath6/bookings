package main

import "testing"

// to run all tests
// go test -v ./...
func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Error("failed run()")
	}
}