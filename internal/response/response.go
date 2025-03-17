package response

import (
    "io"
    "errors"
    "fmt"
    "http-from-tcp/internal/headers"
)

type StatusCode int
const (
    STATUS_OK StatusCode = 200
    STATUS_BAD_REQUEST StatusCode = 400
    STATUS_INTERNAL_SERVER_ERROR StatusCode = 500
)

type WriterState int
const (
    WRITER_STATE_STATUSLINE WriterState = iota
    WRITER_STATE_HEADERS
    WRITER_STATE_BODY
    WRITER_STATE_DONE
)

type Writer struct {
    internalWriter io.Writer
    state WriterState
}

func NewWriter(internalWriter io.Writer) Writer {
    return Writer { internalWriter, WRITER_STATE_STATUSLINE }
}

func (w *Writer) Write(data []byte) (int, error) {
    return w.internalWriter.Write(data)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) (int, error) {
    var statusLine string

    if w.state != WRITER_STATE_STATUSLINE { return 0, errors.New("out of order write") }

    switch(statusCode) {
    case STATUS_OK:
        statusLine = "HTTP/1.1 200 OK"
    case STATUS_BAD_REQUEST:
        statusLine = "HTTP/1.1 400 Bad Request"
    case STATUS_INTERNAL_SERVER_ERROR:
        statusLine = "HTTP/1.1 500 Internal Server Error"
    default:
        return 0, errors.New("unrecognized status code")
    }

    var totalBytes int
    n, err := w.Write([]byte(statusLine))
    totalBytes += n
    if err != nil { return totalBytes, err }
    n, err = w.Write([]byte("\r\n"))
    totalBytes += n
    if err != nil { return totalBytes, err }

    w.state = WRITER_STATE_HEADERS
    return totalBytes, nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) (int, error) {
    if w.state != WRITER_STATE_HEADERS { return 0, errors.New("out of order write") }

    var totalBytes int

    for key, value := range headers {
        header := fmt.Appendf([]byte{}, "%s: %s\r\n", key, value)
        n, err := w.Write(header)
        totalBytes += n
        if err != nil { return totalBytes, err }
    }

    n, err := w.Write([]byte("\r\n"))
    totalBytes += n
    if err != nil { return totalBytes, err }

    w.state = WRITER_STATE_BODY
    return totalBytes, nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
    if w.state != WRITER_STATE_BODY { return 0, errors.New("out of order write") }

    n, err := w.Write(body)
    if err != nil { return n, err }

    w.state = WRITER_STATE_DONE
    return n, nil
}
