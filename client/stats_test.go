package client

import (
	"testing"
	"encoding/json"
	"github.com/sanitizer/test-rest-api/client/mdl"
	"github.com/sanitizer/test-rest-api/client/tls"
	"bytes"
	"net/http"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

/*
	This test was designed to verify data correctness in stats service
	I validate the number of requests to the has service (GET, POST)
	I validate the correctness of average time computations by reverse math
 */
func TestStatsAverageTime(t *testing.T) {
	timeTracking := make([]float64, 0)

	//GET Initial Stats
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	stats := new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	tls.AssertError(e, t)
	initialNumOfReq := stats.Requests
	initialTime := stats.Time
	initialSum := getSumFromAverage(initialTime, initialNumOfReq)

	// SETUP --------------------------------------------------------------------------
	req, _ := json.Marshal(mdl.ReqBody{Password:"t"})
	start := time.Now()
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	end := time.Now()
	timeTracking = append(timeTracking, end.Sub(start).Seconds() * 1000)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t)

	req, _ = json.Marshal(mdl.ReqBody{Password:"te"})
	start = time.Now()
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	end = time.Now()
	timeTracking = append(timeTracking, end.Sub(start).Seconds() * 1000)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t)

	req, _ = json.Marshal(mdl.ReqBody{Password:"tes"})
	start = time.Now()
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	end = time.Now()
	timeTracking = append(timeTracking, end.Sub(start).Seconds() * 1000)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t)
	// ------------------------------------------------------------------------------------------------

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	stats = new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	lowerBound, upperBound := getBoundaries(timeTracking)
	finalSum := getSumOfSubSetFromAverage(stats.Time, stats.Requests, initialSum, len(timeTracking))

	tls.AssertEquals(initialNumOfReq + 3, stats.Requests, t)
	tls.AssertWithin(lowerBound, upperBound, finalSum, t)
}

/*
	This test specifically checks requests piece of data provided by /stats endpoint
	Validating correctness of counting requests to the /hash endpoints (GET, POST)
 */
func TestStatsNumOfRequests(t *testing.T) {
	response, e := tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	stats := new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	tls.AssertError(e, t)
	initialNumOfReq := stats.Requests

	req, _ := json.Marshal(mdl.ReqBody{Password:"te"})
	response, e = tls.SendRequest(mdl.POST_METHOD, mdl.HASH, bytes.NewBuffer(req), mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusCreated, response.StatusCode(), t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	stats = new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	tls.AssertError(e, t)
	tls.AssertEquals(initialNumOfReq + 1, stats.Requests, t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.HASH, nil, mdl.JSON)
	tls.AssertError(e, t)
	tls.AssertEquals(http.StatusOK, response.StatusCode(), t)

	response, e = tls.SendRequest(mdl.GET_METHOD, mdl.STATS, nil, mdl.JSON)
	tls.AssertError(e, t)
	stats = new(mdl.Stats)
	e = json.Unmarshal(response.Body(), stats)
	tls.AssertError(e, t)
	tls.AssertEquals(initialNumOfReq + 2, stats.Requests, t)
}

func getBoundaries(allTimes []float64) (lowerBound int, upperBound int) {
	var lowerBoundTotal float64
	var upperBoundTotal float64

	for _, val := range allTimes {
		//substitute average upper bound latency in mills(between 5 - 40 based on google search)
		lowerBoundTotal += val - mdl.AVG_LATENCY
		upperBoundTotal += val
	}

	lowerBound = int(lowerBoundTotal / float64(len(allTimes)))
	upperBound = int(upperBoundTotal / float64(len(allTimes)))
	return
}

func getSumFromAverage(average int, numOfSamples int) int {
	return average * numOfSamples
}

func getSumOfSubSetFromAverage(average int, numOfSamples int, initialSum int, sampleSize int) int {
	return int((getSumFromAverage(average, numOfSamples) - initialSum) / sampleSize)
}

