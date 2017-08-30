package client

import
(
	"testing"
	"fmt"
	"encoding/json"
	"bytes"
	"io"
	"strconv"
	"github.com/go-resty/resty"
	"errors"
)

const (
	POST_METHOD = "POST"
	GET_METHOD = "GET"
	BASE = "http://radiant-gorge-83016.herokuapp.com/"
	STATS = BASE + "stats"
	HASH = BASE + "hash"
)

type Stats struct {
	Requests int `json:"TotalRequests" xml:"TotalRequests" text:"TotalRequests"`
	Time 	 int `json:"AverageTime" xml:"AverageTime" text:"AverageTime"`
}

func (this Stats) String() string {
	return fmt.Sprintf("Stats ---->\n\tTotal Requests:%d\n\tAverage Time:%d", this.Requests, this.Time)
}

type ReqBody struct {
	Password string `json:"password" xml:"password" text:"password"`
}

func (this ReqBody) String() string {
	return fmt.Sprintf("Request body ---->\n\tPassword:%q", this.Password)
}

func TestSmokeGetRequests(t *testing.T) {
	//GET Stats
	response, e := sendRequest(GET_METHOD, STATS, nil)
	assertError(e, t)
	stats := new(Stats)
	e = json.Unmarshal(response.Body(), stats)
	assertError(e, t)
	fmt.Println("body: ", stats.String())
	fmt.Sprintf("response code: %q", strconv.Itoa(response.StatusCode()))
	fmt.Println("response status: ", response.Status())

	//POST Password
	req, _ := json.Marshal(ReqBody{Password:"test"})
	response, e = sendRequest(POST_METHOD, HASH, bytes.NewBuffer(req))
	jobId := string(response.Body())
	assertError(e, t)
	fmt.Println("body: ", jobId)
	fmt.Sprintf("response code: %q", strconv.Itoa(response.StatusCode()))
	fmt.Println("response status: ", response.Status())

	//GET HASH
	response, e = sendRequest(GET_METHOD, HASH + "/" + jobId, nil)
	assertError(e, t)
	hash := string(response.Body())
	fmt.Println("body: ", hash)
	fmt.Sprintf("response code: %q", strconv.Itoa(response.StatusCode()))
	fmt.Println("response status: ", response.Status())
}

func sendRequest(method string, url string, body io.Reader) (*resty.Response, error) {
	switch method {
	case GET_METHOD:
		return resty.R().SetBody(body).Get(url)
	case POST_METHOD:
		return resty.R().SetBody(body).Post(url)
	}
	return nil, errors.New("Was not able to identify the method...")
}


func TestSmokePostRequests(t *testing.T) {

}


func assert(expected string, actual string, t *testing.T, testId string) {
	if expected != actual {
		t.Errorf("Assertion failure in %q. expected: %q\tactual: %q", testId, expected, actual)
	}

	t.Log("----> ", testId, " Passed")
}

func assertError(e error, t *testing.T) {
	if e != nil {
		t.Errorf("Failed based on received error: %q", e.Error())
	}
}