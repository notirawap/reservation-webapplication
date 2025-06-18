package main

import (
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	mux := routes()
	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("type is not *chi.Mux, type is %T", v)
	}
}
