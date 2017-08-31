package tls

import (
	"testing"
	"io"
	"github.com/go-resty/resty"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"errors"
)

func AssertEquals(expected interface{}, actual interface{}, t *testing.T, testId string) {
	if expected != actual {
		switch expected.(type) {
		case int:
			t.Errorf("Assertion failure in %q. expected: %d\tactual: %d", testId, expected.(int), actual.(int))
		case string:
			t.Errorf("Assertion failure in %q. expected: %q\tactual: %q", testId, expected, actual)
		}
	}else {
		t.Log("----> ", testId, " Passed")
	}
}

func AssertNotEquals(expected interface{}, actual interface{}, t *testing.T, testId string) {
	if expected == actual {
		switch expected.(type) {
		case int:
			t.Errorf("Assertion failure in %q. expected: %d\tactual: %d SHOULD NOT EQUAL", testId, expected.(int), actual.(int))
		case string:
			t.Errorf("Assertion failure in %q. expected: %q\tactual: %q SHOULD NOT EQUAL", testId, expected, actual)
		}
	}else {
		t.Log("----> ", testId, " Passed")
	}
}

func AssertWithin(lower int, upper int, value int, t *testing.T, testId string) {
	if value > upper || value < lower {
		t.Errorf("Assertion failure in %q. lower: %d\tupper: %d\tvalue: %d", testId, lower, upper, value)
	}else {
		t.Log("----> ", testId, " Passed")
	}
}

func AssertNotEmpty(value string, testId string, t *testing.T) {
	if value == "" {
		t.Errorf("Assertion failure in %q. Value is empty.", testId)
	} else {
		t.Log("----> ", testId, " Passed")
	}
}

func AssertEmpty(value string, testId string, t *testing.T) {
	if value != "" {
		t.Errorf("Assertion failure in %q. Value is not empty. Value: %q", testId, value)
	} else {
		t.Log("----> ", testId, " Passed")
	}
}

func AssertError(e error, t *testing.T) {
	if e != nil {
		t.Errorf("Failed based on received error: %q", e.Error())
	}
}

func SendRequest(method string, url string, body io.Reader, contentType string) (*resty.Response, error) {
	switch method {
	case mdl.GET_METHOD:
		return resty.R().SetHeader("Content-type", contentType).SetBody(body).Get(url)
	case mdl.POST_METHOD:
		return resty.R().SetBody(body).Post(url)
	}
	return nil, errors.New("Was not able to identify the method...")
}