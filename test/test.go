package test

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type testBody struct {
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
}

func newRequest(method, path string) *http.Request {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		log.Panicf("error at newRequest: %s", err)
	}
	return req
}

func newRequestWithBody(method, path string, data *url.Values) *http.Request {
	encodedData := data.Encode()
	req, err := http.NewRequest(method, path, strings.NewReader(encodedData))
	if err != nil {
		log.Panicf("error at newRequestWithBody: %s", err)
	}
	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	return req
}
