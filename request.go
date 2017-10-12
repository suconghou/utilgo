package utilgo

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// GetResp return http response repleace for http.Get
func GetResp(url string, timeout uint) (*http.Response, error) {
	return Dohttp(url, "GET", nil, nil, timeout, nil)
}

// GetContent send get request and read response
func GetContent(url string, timeout uint) ([]byte, error) {
	resp, err := GetResp(url, timeout)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

// PostContent send post request and read response
func PostContent(url string, contentType string, body io.Reader, callback func(resp *http.Response) ([]byte, error)) ([]byte, error) {
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	if callback == nil {
		defer resp.Body.Close()
		bodyStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return bodyStr, err
		}
		return bodyStr, nil
	}
	return callback(resp)
}

// PostContentWait send post request and wait one second then read response
func PostContentWait(url string, contentType string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	time.Sleep(time.Second)
	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bodyStr, err
	}
	return bodyStr, nil
}

// Dohttp do a http request and return http response
func Dohttp(url string, method string, reqHeader http.Header, body io.Reader, timeout uint, transport *http.Transport) (*http.Response, error) {
	client := NewClient(timeout, transport)
	req, err := NewRequest(url, method, reqHeader, body)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// NewRequest return *http.Request
func NewRequest(url string, method string, reqHeader http.Header, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	for key, value := range reqHeader {
		for _, item := range value {
			req.Header.Set(key, item)
		}
	}
	return req, nil
}

// NewClient return *http.Client
func NewClient(timeout uint, transport *http.Transport) *http.Client {
	var client *http.Client
	if transport != nil {
		client = &http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: transport}
	} else {
		client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	}
	return client
}
