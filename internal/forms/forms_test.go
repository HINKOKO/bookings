package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid, expected to be valid form")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("t", "r", "v")
	if form.Valid() {
		t.Error("Valid form but required field are missing !")
	}

	postedData := url.Values{}
	postedData.Add("t", "t")
	postedData.Add("r", "r")
	postedData.Add("v", "v")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("t", "r", "v")
	if !form.Valid() {
		t.Error("Validity is not valid")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("GET", "/whatever", nil)
	form := New(r.PostForm)
	has := form.Has("whatever")
	if has {
		t.Error("form show has field when it does not")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("Form does not show has field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("GET", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10) // Checking len of non existing field
	if form.Valid() {
		// So should not be valid
		t.Error("Show a length (min) for non-existing field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Expected an error to be found but it's missing")
	}

	postedValues := url.Values{}
	postedValues.Add("key", "michael")
	form = New(postedValues)

	form.MinLength("key", 100)
	if form.Valid() {
		t.Error("Shows min length of 100 met when data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another", "abc123")
	form = New(postedValues)

	form.MinLength("another", 1)
	if !form.Valid() {
		t.Error("Shows min length of 1 is not met when it is")
	}
	isError = form.Errors.Get("another")
	if isError != "" {
		t.Error("Shoudl NOT have an error but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("X")
	if form.Valid() {
		t.Error("Output a valid a form with unproper email")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "johndoe@gmail.com")
	form = New(postedValues)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Does not valid the proper john doe email when it should")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "@x")
	form = New(postedValues)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Got valid email for invalid email address")
	}
}
