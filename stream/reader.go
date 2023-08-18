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
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF && len(line) == 0 { // the 2nd read after first EOF
			if r.buf.Len() == 0 {
				return "", err // stream reaches the end. caller should get EOF error
			}
			break
		}
		if bytes.HasSuffix(line, []byte("\r\n")) {
			r.buf.Write(bytes.TrimRight(line, "\r\n"))
			break
		}
		r.buf.Write(line)
	}
	return r.buf.String(), nil
}
