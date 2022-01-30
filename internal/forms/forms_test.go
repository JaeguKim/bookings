package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a","a")
	form := New(postedData)
	if !form.Has("a") {
		t.Error("field should have the field")
	}

	postedData = url.Values{}
	form = New(postedData)
	if form.Has("a") {
		t.Error("field should not have the field")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a","abcd")
	form := New(postedData)
	form.MinLength("a",4)
	if !form.Valid() {
		t.Error("field length should be more than min length constraint")
	}
	isError := form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")
	form = New(postedData)
	form.MinLength("a",4)
	if form.Valid() {
		t.Error("min length check failed!")
	}

	isError = form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email","a@b.com")
	form := New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("email is right form but returned false")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("email is not right form but returned true")
	}
}

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a","b","c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a","a")
	postedData.Add("b","a")
	postedData.Add("c","a")
	r,_ = http.NewRequest("POST","/whatever",nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a","b","c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}
