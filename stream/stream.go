package stream

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	streamURL  = "https://api.twitter.com/2/tweets/search/stream"
	maxRetries = 5
)

type Streamer interface {
	Stream() (chan string, chan error, chan struct{}, error)
	Close() error
}

type TwitterStream struct {
	done chan struct{}
}

func NewTwitterStream() *TwitterStream {
	return &TwitterStream{
		done: make(chan struct{}),
	}
}

func (t *TwitterStream) Stream() (chan string, chan error, chan struct{}, error) {
	resp, err := connectTwitterStream()
	if err != nil {
		return nil, nil, nil, err
	}
	msgs := make(chan string)
	errs := make(chan error)
	stop := make(chan struct{})
	go readMessages(resp, msgs, errs, t.done, stop)
	return msgs, errs, stop, nil

}

func (t *TwitterStream) Close() {
	close(t.done)
}

func connectTwitterStream() (*http.Response, error) {
	token := os.Getenv("TWITTER_BEARER_TOKEN")
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
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

func readMessages(resp *http.Response, msgs chan<- string, errs chan<- error, done <-chan struct{}, stop chan<- struct{}) {
	defer func() {
		resp.Body.Close()
		close(msgs)
		close(errs)
		close(stop)
	}()
	reader := ResponseBodyReader{
		reader: bufio.NewReader(resp.Body),
	}
	for {
		select {
		case <-done:
			return
		default:
			msg, err := reader.NextMessage()
			if err != nil {
				errs <- err
				return
			}
			msgs <- msg
		}
	}
}
