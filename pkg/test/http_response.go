package test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	resty "gopkg.in/resty.v1"
)

// CreateHTTPResponse - helper to create an http response for test cases
// body
// parameters:
//   response code:
//   data[0] response body
//   data[1] http status
//   data[x*2 - 1] = http header where x > 1
//   data[x*2] = http header value
func CreateHTTPResponse(respcode int, data ...string) *resty.Response {
	resp := http.Response{}

	resp.Header = make(http.Header)
	resp.StatusCode = respcode

	if len(data) < 2 { // if no body is provided, create one
		resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("")))
	}

	for k, v := range data {
		if k == 0 {
			resp.Status = v
			continue
		}
		if k == 1 {
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(v)))
			resp.ContentLength = int64(len(v))

			continue
		}
		if k%2 == 0 {
			continue
		}
		resp.Header.Add(data[k-1], v)
	}
	response := &resty.Response{
		RawResponse: &resp,
	}

	return response
}