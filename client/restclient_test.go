package client

import
(
	"testing"
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"bytes"
	"io"
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
	response, e := sendRequst(GET_METHOD, STATS, nil)
	handleError(e)
	stats := new(Stats)
	e = json.Unmarshal(response, stats)
	handleError(e)
	fmt.Println("body: ", stats.String())

	//POST Password
	req, _ := json.Marshal(ReqBody{Password:"test"})
	response, e = sendRequst(POST_METHOD, HASH, bytes.NewBuffer(req))
	jobId := string(response[:])
	handleError(e)
	fmt.Println("body: ", jobId)

	//GET HASH
	response, e = sendRequst(GET_METHOD, HASH + "/" + jobId, nil)
	handleError(e)
	hash := string(response[:])
	fmt.Println("body: ", hash)
}


func TestSmokePostRequests(t *testing.T) {

}

func sendRequst(method string, url string, body io.Reader) ([]byte, error)  {
	request, e := http.NewRequest(method, url, body)

	if e != nil {
		return nil, e
	}

	request.Header.Set("Content-type", "application/json")
	cl := new (http.Client)
	response, e := cl.Do(request)
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func handleError(e error) {
	if e != nil {
		log.Fatal("Error reading: " + e.Error())
	}
}