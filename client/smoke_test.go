package client

import (
	"testing"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"github.com/sanitizer/test-rest-api/client/tls"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
)

func TestStatsGetMethodStatusOK(t *testing.T) {
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t, "Test Stats endpoint without request body. Status OK")

	reqBody := []byte(`{"hello": "world"}`)
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t, "Test Stats endpoint with request body. Status OK")
}

func TestHashPostMethodStatusOk(t *testing.T) {
	req, _ := json.Marshal(mdl.ReqBody{Password:"test"})
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t, "Test Hash endpoint with request body. Status OK")
	jobId := string(response.Body())
	tls.AssertNotEmpty(jobId, "Test Hash endpoint body returned job id", t)
}

func TestHashGetMethodStatusOk(t *testing.T) {
	//Precondition (create a job id first)
	req, _ := json.Marshal(mdl.ReqBody{Password:"test"})
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	jobId := string(response.Body())
	tls.AssertNotEmpty(jobId, "Test Hash endpoint body returned job id", t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t, "Test Hash endpoint with known job id. Status OK")
	hash := string(response.Body())
	tls.AssertNotEmpty(hash,"Test Hash endpoint with known job id. Status OK", t)
	tls.AssertEquals(mdl.HASH_LEN, len(hash),t, "Test Hash endpoint with known job id. Hash length test")

	//Precondition (create a job id first)
	req, _ = json.Marshal(mdl.ReqBody{Password:"ThisistheTest"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	jobId = string(response.Body())
	tls.AssertNotEmpty(jobId, "Test Hash endpoint body returned job id", t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t, "Test Hash endpoint with known job id. Status OK")
	hash2 := string(response.Body())
	tls.AssertNotEmpty(hash,"Test Hash endpoint with known job id. Status OK", t)
	tls.AssertEquals(mdl.HASH_LEN, len(hash2),t, "Test Hash endpoint with known job id. Hash length test")
	tls.AssertNotEquals(hash, hash2,t, "Test Hash endpoint with known job id. Hashs non equality test")

	//Precondition (create a job id first)
	req, _ = json.Marshal(mdl.ReqBody{Password:"ThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTest"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	jobId = string(response.Body())
	tls.AssertNotEmpty(jobId, "Test Hash endpoint body returned job id", t)

	reqBody := []byte(`{"hello": "world"}`)
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t, "Test Hash endpoint with known job id with body. Status OK")
	hash = string(response.Body())
	tls.AssertNotEmpty(hash,"Test Hash endpoint with known job id. Status OK", t)
	tls.AssertEquals(mdl.HASH_LEN, len(hash),t, "Test Hash endpoint with known job id. Hash length test")
}

func TestHashGetMethodStatusNotOk(t *testing.T) {
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.HASH, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusMethodNotAllowed, response.StatusCode(), t, "Test Hash endpoint without job id. Status Method not allowed")

	reqBody := []byte(`{"hello": "world"}`)
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusMethodNotAllowed, response.StatusCode(), t, "Test Hash endpoint without job id with body. Status Method not allowed")

	//request with job id that is not a number
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/test", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t, "Test Hash endpoint with job id that is NAN. Status Bad Request")

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/-1", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t, "Test Hash endpoint with job id that does not exist. Status Bad Request")

	//The assertions bellow assume that correct response would be 404 and not 400 as it is now
	//request with job id of bigger than max uint64 in golang
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + strconv.FormatUint(mdl.MAX_UINT, 10) + "6", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusNotFound, response.StatusCode(), t, "Test Hash endpoint with job id that is higher than int in go. Status Not Found")

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/6666666", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusNotFound, response.StatusCode(), t, "Test Hash endpoint with job id that is just a high number. This test can be false negative. Status Not Found")
}

func TestHashPostMethodWithEmptyBodyOrEmptyPasswStatusNotOk(t *testing.T) {
	req, _ := json.Marshal(mdl.ReqBody{Password:""})
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t, "Test Hash endpoint with request body but empty password. Status Bad request")
	jobId := string(response.Body())
	tls.AssertEmpty(jobId, "Test Hash endpoint body returned job id", t)

	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t, "Test Hash endpoint without request body. Status Bad request")
	jobId = string(response.Body())
	tls.AssertEquals("Malformed Input\n", jobId, t, "Test Hash endpoint without request body returns error msg in body")

	req, _ = json.Marshal(mdl.ReqBody{Password:"0 or 1=1"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t, "Test Hash endpoint without SQL Injection 1. Status Bad request")
	jobId = string(response.Body())
	tls.AssertEquals("Malformed Input\n", jobId, t, "Test Hash endpoint without SQL Injection 1 with error in body")

	req, _ = json.Marshal(mdl.ReqBody{Password:"0; Drop user administrator"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t, "Test Hash endpoint without SQL Injection 2. Status Bad request")
	jobId = string(response.Body())
	tls.AssertEquals("Malformed Input\n", jobId, t, "Test Hash endpoint without SQL Injection 2 with error in body")
}

// this is commented for now as it takes forever to restart the service
//func TestHashShutdownPostMethodStatusOkRequestAfterStatusNoResponse(t *testing.T) {
//	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(mdl.SHUTDOWN_REQ_BODY), mdl.JSON)
//	tls.AssertError(e, t)
//	tls.AssertEquals(http.StatusOK, response.StatusCode(), t, "Test Hash Shutdown. Status OK")
//
//	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
//	tls.AssertError(e, t)
//	tls.AssertEquals(http.StatusServiceUnavailable, response.StatusCode(), t, "Test Stats after shutdown was initiated. Status Service Unavailable")
//}