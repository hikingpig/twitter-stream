package request

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	maxRetries = 5
)

func Request(method string, url string, body io.Reader) (*http.Response, error) {
	token := os.Getenv("TWITTER_BEARER_TOKEN")
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	resp, err = checkConnectionStatus(resp, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func checkConnectionStatus(resp *http.Response, req *http.Request) (*http.Response, error) {
	var err error
	if resp.StatusCode == 429 {
		resp, err = backoffRetry(req)
		if err != nil {
			return nil, err
		}
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to connect %s with status code %d", req.URL.String(), resp.StatusCode)
	}
	return resp, nil
}
func backoffRetry(req *http.Request) (*http.Response, error) {
	for i := 0; i < maxRetries; i++ {
		time.Sleep(getBackoffTime(i))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 429 {
			return resp, err
		}
	}
	return nil, errors.New("max retries exceeded")
}

func getBackoffTime(retryNum int) time.Duration {
	n := 1
	for i := 0; i < retryNum; i++ {
		n *= 2
	}
	return time.Duration(n) * time.Second
}
