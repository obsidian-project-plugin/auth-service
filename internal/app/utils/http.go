package utils

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewHTTPRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func PostForm(url string, data url.Values) (*http.Response, error) {
	reqBody := strings.NewReader(data.Encode())
	req, err := NewHTTPRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
