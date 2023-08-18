package stream

import (
	"net/http"
	"testing"
)

func TestStop(t *testing.T) {
	stream := &TwitterStream{
		done: make(chan struct{}),
	}
	stream.Close()
	_, ok := <-stream.done
	if ok != false {
		t.Errorf("expected channel closed, got ok = %v", ok)
	}
}

type mockBodyResponseReader struct{}

func (r *mockBodyResponseReader) SetBody(*http.Response) {}

func (r *mockBodyResponseReader) NextMessage() (string, error) {
	return "hello", nil
}

func mockConnectTwitter() (*http.Response, error) {
	return &http.Response{}, nil
}

func TestStream(t *testing.T) {
	stream := &TwitterStream{
		done:           make(chan struct{}),
		connectTwitter: mockConnectTwitter,
		reader:         &mockBodyResponseReader{},
	}
	msgs, _, _, err := stream.Stream()
	if err != nil {
		t.Errorf("got err when starting stream %v", err)
	}
	msg := <-msgs
	expected := "hello"
	if msg != expected {
		t.Errorf("got %s, want %s", msg, expected)
	}
}
