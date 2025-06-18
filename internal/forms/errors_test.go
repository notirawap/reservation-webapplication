package forms

import "testing"

var e errors

func TestForms_Add(t *testing.T) {
	e.Add("error", "error message")
}

func TestForms_Get(t *testing.T) {
	resp := e.Get("error")
	if resp != "" {
		t.Error("Get the value but should not have got the value")
	}

	e.Add("error", "error message")
	resp = e.Get("error")
	if resp == "error message" {
		t.Error("Cannot get value but should have got the value")
	}
}
