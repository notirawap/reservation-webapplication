package forms

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Server Side Form validation - If error, populate errors field
// Check for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Check for length and type integer
func (f *Form) NumLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) != length {
		f.Errors.Add(field, fmt.Sprintf("This field must be %d characters long", length))
		return false
	} else if _, err := strconv.Atoi(x); err != nil {
		f.Errors.Add(field, "This field must be a number")
		return false
	}
	return true
}

// Check for valid email address using govalidator thrid party package
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
