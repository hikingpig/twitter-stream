package stream

import "testing"

func TestStopStream(t *testing.T) {
	stream := &TwitterStream{
		done: make(chan struct{}),
	}
	stream.Close()
	_, ok := <-stream.done
	if ok != false {
		t.Errorf("expected channel closed, got ok = %v", ok)
	}
}
