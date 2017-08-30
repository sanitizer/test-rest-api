package client

import
(
	"testing"
	"fmt"
	"encoding/json"
	"bytes"
	"strconv"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"github.com/sanitizer/test-rest-api/client/tls"
)

func TestSmokeGetRequests(t *testing.T) {

	//GET Stats
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	stats := new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	tls.AssertError(e, t)
	fmt.Println("body: ", stats.String())
	fmt.Sprintf("response code: %q", strconv.Itoa(response.StatusCode()))
	fmt.Println("response status: ", response.Status())

	//POST Password
	req, _ := json.Marshal(mdl.ReqBody{Password:"test"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	jobId := string(response.Body())
	tls.AssertError(e, t)
	fmt.Println("body: ", jobId)
	fmt.Sprintf("response code: %q", strconv.Itoa(response.StatusCode()))
	fmt.Println("response status: ", response.Status())

	//GET HASH
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, nil, mdl.JSON)
	tls.AssertError(e, t)
	hash := string(response.Body())
	fmt.Println("body: ", hash)
	fmt.Sprintf("response code: %q", strconv.Itoa(response.StatusCode()))
	fmt.Println("response status: ", response.Status())
}

func TestSmokePostRequests(t *testing.T) {

}


