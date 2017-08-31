package client

import (
	"testing"
	"github.com/sanitizer/test-rest-api/client/tls"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func TestHashPostMethodWith1000Requests(t *testing.T) {
	var counter int
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req, e := json.Marshal(mdl.ReqBody{Password: "tesdgsdgsetsgdsdg"})
			tls.AssertError(e, t)
			log.Println("Running routine #", i)
			response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
			if http.StatusCreated != response.StatusCode() {
				counter++
			}
			tls.AssertEquals(http.StatusCreated, response.StatusCode(), t, "Creating hash with routine #" + strconv.Itoa(i))
		}(i)
	}

	wg.Wait()
	tls.AssertEquals(0, counter, t, "Validate number of failed routines in POST")
}

func TestHashGetMethodWith1000Requests(t *testing.T) {
	var counter int
	var wg sync.WaitGroup

	req, e := json.Marshal(mdl.ReqBody{Password: "tesdgsdgsetsgdsdg"})
	tls.AssertError(e, t)
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t, "Creating initial hash")
	jobId := string(response.Body())

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(i int, jbId string) {
			defer wg.Done()
			log.Println("Running routine #", i)
			r, e := tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, nil, mdl.JSON)
			tls.AssertError(e, t)

			if http.StatusOK != r.StatusCode() {
				counter++
			}
			tls.AssertEquals(http.StatusOK, r.StatusCode(), t, "Getting hash with routine #" + strconv.Itoa(i))
		}(i, jobId)
	}

	wg.Wait()
	tls.AssertEquals(0, counter, t, "Validate number of failed routines in GET")
}