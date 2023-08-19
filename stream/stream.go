package stream

import (
	"net/http"

	"github.com/hikingpig/twitter-stream/request"
)

const (
	streamURL = "https://api.twitter.com/2/tweets/search/stream"
)

type Streamer interface {
	Stream() (chan string, chan error, chan struct{}, error)
	Close()
}

type TwitterStream struct {
	done           chan struct{}
	connectTwitter func() (*http.Response, error)
	reader         IResponseBodyReader
}

func NewTwitterStream() *TwitterStream {
	return &TwitterStream{
		done:           make(chan struct{}),
		connectTwitter: connectTwitterStream,
		reader:         &ResponseBodyReader{},
	}
}

func (t *TwitterStream) Stream() (chan string, chan error, chan struct{}, error) {
	resp, err := t.connectTwitter()
	if err != nil {
		return nil, nil, nil, err
	}
	msgs := make(chan string)
	errs := make(chan error)
	stop := make(chan struct{})
	go t.readMessages(resp, msgs, errs, stop)
	return msgs, errs, stop, nil
}

func (t *TwitterStream) Close() {
	close(t.done)
}

func connectTwitterStream() (*http.Response, error) {
	return request.Request(http.MethodGet, streamURL, nil)
}

func (t *TwitterStream) readMessages(resp *http.Response, msgs chan<- string, errs chan<- error, stop chan<- struct{}) {
	defer func() {
		resp.Body.Close()
		close(msgs)
		close(errs)
		close(stop)
	}()
	t.reader.SetBody(resp)
	for {
		select {
		case <-t.done:
			return
		default:
			msg, err := t.reader.NextMessage()
			if err != nil {
				errs <- err
				return
			}
			msgs <- msg
		}
	}
}
