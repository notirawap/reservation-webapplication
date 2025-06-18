package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	// Create form struct (argument)
	postedData := url.Values{}
	form := New(postedData)

	// test valid response
	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Has(t *testing.T) {
	// Create form struct (argument)
	postedData := url.Values{}
	form := New(postedData)

	// test invalid response
	has := form.Has("a")
	if has {
		t.Error("form shows has field when the field is missing")
	}

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	// test valid response
	has = form.Has("a")
	if !has {
		t.Error("form shows has no field when the field is presented")
	}

}

func TestForm_Required(t *testing.T) {
	// Create form struct (argument)
	postedData := url.Values{}
	form := New(postedData)

	// test invalid response
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")
	form = New(postedData)

	// test valid response
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows missing fields but all the required fields present")
	}
}

func TestForm_NumLength(t *testing.T) {
	// Create form struct (argument)
	postedData := url.Values{}
	form := New(postedData)

	// test missing field
	form.NumLength("a", 3)
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("a", "aa")
	form = New(postedData)

	// test invalid response (incorrect length and type)
	length := 5
	form.NumLength("a", length)
	if form.Valid() {
		t.Errorf("form shows valid fields but all the fields' length are not equal to %d and not integer", length)
	}

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("a", "22")
	form = New(postedData)

	// test invalid response (incorrect length)
	form.NumLength("a", length)
	if form.Valid() {
		t.Errorf("form shows valid fields but all the fields' length are not equal to %d", length)
	}

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("a", "aaaaa")
	form = New(postedData)

	// test invalid response (incorrect type)
	form.NumLength("a", length)
	if form.Valid() {
		t.Errorf("form shows valid fields but all the fields' length are not integer")
	}

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("a", "22222")
	form = New(postedData)

	// test valid response
	form.NumLength("a", length)
	if !form.Valid() {
		t.Error("form shows invalid fields but all the required fields are valid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// Create form struct (argument)
	postedData := url.Values{}
	form := New(postedData)

	// test missing data
	form.IsEmail("a")
	if form.Valid() {
		t.Error("form shows valid when required field is missing")
	}

	// ##########Test Get() method of errors struct##########
	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("Should have had an error!")
	}
	// #####################################################

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	// test invalid response
	form.IsEmail("a")
	if form.Valid() {
		t.Error("form shows valid when required fields are invalid")
	}

	// Create form struct (argument)
	postedData = url.Values{}
	postedData.Add("email", "a@test.com")
	form = New(postedData)

	// test valid response
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("form shows invalid email address when the email address field is valid")
	}

	// ##########Test Get() method of errors struct##########
	isError = form.Errors.Get("email")
	if isError != "" {
		t.Error("Should not have had an error!")
	}
	// #####################################################
}
