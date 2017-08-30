package client

import (
	"testing"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"github.com/sanitizer/test-rest-api/client/tls"
	"net/http"
	"bytes"
)

func TestStatsGetMethodStatusOK(t *testing.T) {
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.Assert(http.StatusOK, response.StatusCode(), t, "Test Stats endpoint without request body. Status OK")

	reqBody := []byte(`{"hello": "world"}`)
	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.Assert(http.StatusOK, response.StatusCode(), t, "Test Stats endpoint with request body. Status OK")
}

func TestStatsPostMethodStatusNotOk(t *testing.T) {
	reqBody := []byte(`{"hello": "world"}`)
	response, e := tls.SendRequest(mdl.POST_METHOD, mdl.STATS, bytes.NewBuffer(reqBody), mdl.JSON)
	tls.AssertError(e, t)
	tls.Assert(http.StatusNotFound, response.StatusCode(), t, "Test Stats endpoint with POST method. Status Not OK")
}