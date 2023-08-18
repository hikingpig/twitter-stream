package stream

import (
	"bufio"
	"bytes"
	"io"
)

type ResponseBodyReader struct {
	reader *bufio.Reader
	buf    bytes.Buffer
}

func (r *ResponseBodyReader) NextMessage() (string, error) {
	r.buf.Truncate(0)
	for {
		line, err := r.reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		if bytes.HasSuffix(line, []byte("\r\n")) {
			r.buf.Write(bytes.TrimRight(line, "\r\n"))
			break
		}
		r.buf.Write(line)
	}
	return r.buf.String(), nil
}
