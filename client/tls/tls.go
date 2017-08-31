package tls

import (
	"testing"
	"io"
	"github.com/go-resty/resty"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"errors"
	"runtime/debug"
)

func AssertEquals(expected interface{}, actual interface{}, t *testing.T) {
	if expected != actual {
		switch expected.(type) {
		case int:
			debug.PrintStack()
			t.Errorf("Equals Assertion failure. Expected: %d\tActual: %d", expected.(int), actual.(int))
		case string:
			debug.PrintStack()
			t.Errorf("Equals Assertion failure. Expected: %d\tActual: %d", expected, actual)
		}
	}
}

func AssertNotEquals(expected interface{}, actual interface{}, t *testing.T) {
	if expected == actual {
		switch expected.(type) {
		case int:
			debug.PrintStack()
			t.Errorf("Not Equal Assertion failure. Expected: %d\tActual: %d", expected.(int), actual.(int))
		case string:
			debug.PrintStack()
			t.Errorf("Not Equal Assertion failure. Expected: %q\tActual: %q", expected, actual)
		}
	}
}

func AssertWithin(lower int, upper int, value int, t *testing.T) {
	if value > upper || value < lower {
		debug.PrintStack()
		t.Errorf("Within Assertion failure. Lower: %d\tUpper: %d\tValue: %d", lower, upper, value)
	}
}

func AssertNotEmpty(value string, t *testing.T) {
	if value == "" {
		debug.PrintStack()
		t.Errorf("Not Empty Assertion failure.")
	}
}

func AssertEmpty(value string, t *testing.T) {
	if value != "" {
		debug.PrintStack()
		t.Errorf("Empty Assertion failure. Value: %q", value)
	}
}

func AssertError(e error, t *testing.T) {
	if e != nil {
		debug.PrintStack()
		t.Errorf("Error Assertion failure: %q", e.Error())
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