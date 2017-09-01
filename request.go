package utilgo

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// GetResp return http response
func GetResp(url string) (*http.Response, error) {
	return http.Get(url)
}

// GetContent send get request and read response
func GetContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
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
func PostContent(url string, contentType string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bodyStr, err
	}
	return bodyStr, nil
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
