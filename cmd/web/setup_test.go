package main

// This set up for testing environment that will run before the actual test run

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run()) // Before it exits, it runs test
}

// Create Handler object that satisfies http.Handler interface
type myHandler struct{}
func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
