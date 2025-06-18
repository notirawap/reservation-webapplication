package main

import (

	"net/http"
	"testing"
)

// Testing middleware functions that should receive and return http.Handler

// Test NoSurf
func TestNoSurve(t *testing.T) {
	// Create handler for NoSurf to receive
	var myH myHandler
	h := NoSurf(&myH)

	// NoSurf should return http.Handler
	switch v := h.(type) {
	case http.Handler:
		// pass the test and do nothing
	default:
		t.Errorf("type is not http.Handler, but it is %T", v)
	}
}

// Test SessionLoad
func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
	default:
		t.Errorf("type is not http.Handler, but it is %T", v)
	}
}
