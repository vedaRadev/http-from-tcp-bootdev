package request

import (
    "testing"
    "fmt"
    "io"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type chunkReader struct {
    data string
    numBytesPerRead int
    pos int
}

// Read up to len(p) or numBytesPerRead bytes from the string per call.
// Useful for simulating reading a variaable number of bytes per chunk from
// a network connection.
func (cr *chunkReader) Read(p []byte) (n int, err error) {
    if cr.pos >= len(cr.data) {
        return 0, io.EOF
    }
    endIndex := min(cr.pos + cr.numBytesPerRead, len(cr.data))
    n = copy(p, cr.data[cr.pos:endIndex])
    cr.pos += n
    if n > cr.numBytesPerRead {
        n = cr.numBytesPerRead
        cr.pos -= n - cr.numBytesPerRead
    }
    return n, nil
}

func TestRequestLineParse(t *testing.T) {
    var requestString string
    var reader chunkReader

    // Test: Good GET Request line
    requestString = "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
    fmt.Printf("requestString size: %v\n", len(requestString))
    reader = chunkReader { data: requestString, numBytesPerRead: 1 }
    r, err := RequestFromReader(&reader)
    require.NoError(t, err)
    require.NotNil(t, r)
    assert.Equal(t, "GET", r.RequestLine.Method)
    assert.Equal(t, "/", r.RequestLine.RequestTarget)
    assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

    // Test: Good GET Request line with path
    requestString = "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
    reader = chunkReader { data: requestString, numBytesPerRead: len(requestString) }
    r, err = RequestFromReader(&reader)
    require.NoError(t, err)
    require.NotNil(t, r)
    assert.Equal(t, "GET", r.RequestLine.Method)
    assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
    assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

    // Test: Good POST Request with path
    requestString = "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
    reader = chunkReader { data: requestString, numBytesPerRead: 5 }
    r, err = RequestFromReader(&reader)
    require.NoError(t, err)
    require.NotNil(t, r)
    assert.Equal(t, "POST", r.RequestLine.Method)
    assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
    assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

    // Test: Invalid number of parts in request line
    requestString = "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
    reader = chunkReader { data: requestString, numBytesPerRead: 8 }
    _, err = RequestFromReader(&reader)
    require.Error(t, err)

    // Test: Invalid method (out of order) Request line
    requestString = "/coffee POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
    reader = chunkReader { data: requestString, numBytesPerRead: 3 }
    _, err = RequestFromReader(&reader)
    require.Error(t, err)

    // Test: Invalid version in Request line
    requestString = "OPTIONS /prime/rib TCP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
    reader = chunkReader { data: requestString, numBytesPerRead: 12 }
    _, err = RequestFromReader(&reader)
    require.Error(t, err)
}
