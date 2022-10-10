package utils

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)


func BuildRequestURL(baseUrl string, endpoint string, params map[string]string) (*url.URL, error) {
	reqUrl, err := url.Parse(baseUrl)
	if(err != nil) {
		return nil, err
	}

	reqUrl = reqUrl.JoinPath(endpoint)

	// add the url params 
	if (params != nil) {
		reqParams := reqUrl.Query()
		for paramName, paramValue := range params {
			reqParams.Add(paramName, paramValue)
		}
		reqUrl.RawQuery = reqParams.Encode()
	}
	
	return reqUrl, nil	
}

func sendHttpRequest(method string, URI string, headers map[string]string, postParams map[string]string) (*http.Response, error) {
	postValues := url.Values{}
	for paramName, paramValue := range postParams {
		postValues.Add(paramName, paramValue)
	}

	req, _ := http.NewRequest(method, URI, strings.NewReader(postValues.Encode()))

	for headerName, headerValue := range headers {
		req.Header.Add(headerName, headerValue)
	}
	
	return http.DefaultClient.Do(req)	
}

func SendGetApiRequest(baseUrl string, endpoint string, params map[string]string, headers map[string]string) (*http.Response, error) {

	reqUrl, err := BuildRequestURL(baseUrl, endpoint, params)
	if(err != nil) {
		return nil, errors.Wrap(err, "failed building url")
	}

	return sendHttpRequest(http.MethodGet, reqUrl.String(), headers, nil)
}


func SendPostApiRequest(baseUrl string, endpoint string, params map[string]string, headers map[string]string) (*http.Response, error) {

	reqUrl, err := BuildRequestURL(baseUrl, endpoint, nil)
	if(err != nil) {
		return nil, errors.Wrap(err, "failed building url")
	}

	return sendHttpRequest(http.MethodPost, reqUrl.String(), headers, params)
}

