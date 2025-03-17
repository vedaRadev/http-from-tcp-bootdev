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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
    var statusLine string

    if w.state != WRITER_STATE_STATUSLINE { return errors.New("out of order write") }

    switch(statusCode) {
    case STATUS_OK:
        statusLine = "HTTP/1.1 200 OK"
    case STATUS_BAD_REQUEST:
        statusLine = "HTTP/1.1 400 Bad Request"
    case STATUS_INTERNAL_SERVER_ERROR:
        statusLine = "HTTP/1.1 500 Internal Server Error"
    default:
        return errors.New("unrecognized status code")
    }

    _, err := w.Write([]byte(statusLine))
    if err != nil { return err }
    _, err = w.Write([]byte("\r\n"))
    if err != nil { return err }

    w.state = WRITER_STATE_HEADERS
    return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
    if w.state != WRITER_STATE_HEADERS { return errors.New("out of order write") }

    for key, value := range headers {
        header := fmt.Appendf([]byte{}, "%s: %s\r\n", key, value)
        _, err := w.Write(header)
        if err != nil { return err }
    }

    _, err := w.Write([]byte("\r\n"))
    if err != nil { return err }

    w.state = WRITER_STATE_BODY
    return nil
}

func (w *Writer) WriteBody(body []byte) error {
    if w.state != WRITER_STATE_BODY { return errors.New("out of order write") }

    _, err := w.Write(body)
    if err != nil { return err }

    w.state = WRITER_STATE_DONE
    return nil
}
