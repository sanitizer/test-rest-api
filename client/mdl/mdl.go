package mdl

import "fmt"

const (
	POST_METHOD = "POST"
	GET_METHOD = "GET"
	BASE = "http://radiant-gorge-83016.herokuapp.com/"
	STATS = BASE + "stats"
	HASH = BASE + "hash"
	JSON = "application/json"
	TEXT = "application/text"
	XML = "application/xml"
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