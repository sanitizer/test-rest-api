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

/*
	Testing the correctness of status code
	Testing the fact that although the method does not accept any body with request, the response data should not be affected
 */
func TestStatsGetMethodStatusOK(t *testing.T) {
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t)
	stats := new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	tls.AssertError(e,t)

	reqBody := []byte(`{"hello": "world"}`)
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t)
	stats2 := new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	tls.AssertError(e, t)

	tls.AssertEquals(stats.Requests, stats2.Requests, t)
	tls.AssertEquals(stats.Time, stats2.Time, t)
}

/*
	Testing the correctness of status code
	Testing idempotency of post method
 */
func TestHashPostMethodStatusOk(t *testing.T) {
	req, _ := json.Marshal(mdl.ReqBody{Password:"test"})
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t)
	jobId := string(response.Body())
	tls.AssertNotEmpty(jobId, t)

	//create a hash for the same password
	req, _ = json.Marshal(mdl.ReqBody{Password:"test"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t)
	jobId2 := string(response.Body())
	tls.AssertNotEmpty(jobId, t)
	tls.AssertEquals(jobId, jobId2, t)
}

/*
	Testing the correctness of status code
	Testing retrieval of resources known and unknown
 */
func TestHashGetMethodStatusOk(t *testing.T) {
	//Precondition (create a job id first)
	req, _ := json.Marshal(mdl.ReqBody{Password:"test"})
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	jobId := string(response.Body())
	tls.AssertNotEmpty(jobId, t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t)
	hash := string(response.Body())
	tls.AssertNotEmpty(hash, t)
	tls.AssertEquals(mdl.HASH_LEN, len(hash), t)

	//Precondition (create a job id first)
	req, _ = json.Marshal(mdl.ReqBody{Password:"ThisistheTest"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	jobId = string(response.Body())
	tls.AssertNotEmpty(jobId, t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t)
	hash2 := string(response.Body())
	tls.AssertNotEmpty(hash, t)
	tls.AssertEquals(mdl.HASH_LEN, len(hash2),t)
	tls.AssertNotEquals(hash, hash2, t)

	//Precondition (create a job id first)
	req, _ = json.Marshal(mdl.ReqBody{Password:"ThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTestThisistheTest"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	jobId = string(response.Body())
	tls.AssertNotEmpty(jobId, t)

	reqBody := []byte(`{"hello": "world"}`)
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + jobId, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t)
	hash = string(response.Body())
	tls.AssertNotEmpty(hash, t)
	tls.AssertEquals(mdl.HASH_LEN, len(hash), t)
}

/*
	Testing the correctess of status codes
	Testing numberic resource requests, negative and positive, existing resource and not existing resource
	Testing with no jobId in the resource request
	Testing with non numeric jobId in the resource request
	Testing the fact that although the body in the request is not accepted, the returned data should not be affected
 */
func TestHashGetMethodStatusNotOk(t *testing.T) {
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.HASH, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusMethodNotAllowed, response.StatusCode(), t)

	reqBody := []byte(`{"hello": "world"}`)
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusMethodNotAllowed, response.StatusCode(), t)

	//request with job id that is not a number
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/test", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/-1", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t)

	//The assertions bellow assume that correct response would be 404 and not 400 as it is now
	//request with job id of bigger than max uint64 in golang
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/" + strconv.FormatUint(mdl.MAX_UINT, 10) + "6", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusNotFound, response.StatusCode(), t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/6666666", nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusNotFound, response.StatusCode(), t)
}

/*
	Testing the correctness of status codes
	Testing with sql injections as a password, password that is empty, password with spaces, regular password, empty body with request
 */
func TestHashPostMethodWithEmptyBodyOrEmptyPasswStatusNotOk(t *testing.T) {
	req, _ := json.Marshal(mdl.ReqBody{Password:""})
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t)
	jobId := string(response.Body())
	tls.AssertEmpty(jobId, t)
	if _, e := strconv.Atoi(jobId); e == nil {
		t.Error("Returned jobId when expected error message. JobId: ", jobId)
	}

	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t)
	jobId = string(response.Body())
	if _, e := strconv.Atoi(jobId); e == nil {
		t.Error("Returned jobId when expected error message. JobId: ", jobId)
	}

	req, _ = json.Marshal(mdl.ReqBody{Password:"0 or 1=1"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t)
	jobId = string(response.Body())
	if _, e := strconv.Atoi(jobId); e == nil {
		t.Error("Returned jobId when expected error message. JobId: ", jobId)
	}

	req, _ = json.Marshal(mdl.ReqBody{Password:"0; Drop user administrator"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusBadRequest, response.StatusCode(), t)
	jobId = string(response.Body())
	if _, e := strconv.Atoi(jobId); e == nil {
		t.Error("Returned jobId when expected error message. JobId: ", jobId)
	}
}

/* 	this is commented for now as it takes forever to restart the service
	Testing the correctness of status codes
	Testing graceful shutdown
	Testing the fact that no more requests should be processed while in the process of shutdown
 */
//func TestHashShutdownPostMethodStatusOkRequestAfterStatusNoResponse(t *testing.T) {
//	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(mdl.SHUTDOWN_REQ_BODY), mdl.JSON)
//	tls.AssertError(e, t)
//	tls.AssertEquals(http.StatusOK, response.StatusCode(), t)
//
//	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
//	tls.AssertError(e, t)
//	tls.AssertEquals(http.StatusServiceUnavailable, response.StatusCode(), t)

//	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH + "/1", nil, mdl.JSON)
//	tls.AssertError(e, t)
//	tls.AssertEquals(http.StatusServiceUnavailable, response.StatusCode(), t)

//	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, nil, mdl.JSON)
//	tls.AssertError(e, t)
//	tls.AssertEquals(http.StatusServiceUnavailable, response.StatusCode(), t)
//}