package tls

import (
	"testing"
	"io"
	"github.com/go-resty/resty"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"errors"
	"strconv"
)

func Assert(expected interface{}, actual interface{}, t *testing.T, testId string) {
	if expected != actual {
		switch expected.(type) {
		case int:
			t.Errorf("Assertion failure in %q. expected: %q\tactual: %q", testId, strconv.Itoa(expected.(int)), strconv.Itoa(actual.(int)))
		case string:
			t.Errorf("Assertion failure in %q. expected: %q\tactual: %q", testId, expected, actual)
		}
	}

	t.Log("----> ", testId, " Passed")

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