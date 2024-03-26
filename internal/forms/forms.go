package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form - ccreates custom Form struct, embeds url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid - check if form is valid, all data field satisfyied
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New - init a Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// to perform all check in one go -- Look at the nice variadic parameter 'fields' to range through the fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has - check if a field is present, server side validation
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength - Perform length check on strings in form
func (f *Form) MinLength(field string, length int) bool {
	// x := r.Form.Get(field) // here we woudl check the request associated
	x := f.Get(field) // We should instead check the data associated in the form - for test to work properly
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail - Check for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
