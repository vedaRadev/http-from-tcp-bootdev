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

type Writer struct {
    internalWriter io.Writer
}

func NewWriter(internalWriter io.Writer) Writer {
    return Writer { internalWriter }
}

func (w *Writer) Write(data []byte) (int, error) {
    return w.internalWriter.Write(data)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) (int, error) {
    var statusLine string

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

    return totalBytes, nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) (int, error) {
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

    return totalBytes, nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
    n, err := w.Write(body)
    if err != nil { return n, err }

    return n, nil
}

func (w *Writer) WriteChunkedBody(data []byte) (int, error) {
    var totalBytes int
    dataLen := len(data)

    n, err := w.Write(fmt.Appendf([]byte{}, "%x\r\n", dataLen))
    totalBytes += n
    if err != nil { return totalBytes, err }

    send := make([]byte, dataLen)
    copy(send, data)
    send = append(send, '\r', '\n')
    n, err = w.Write(send)
    totalBytes += n
    return totalBytes, err
}

func (w *Writer) WriteChunkedBodyDone(trailers headers.Headers) (int, error) {
    var totalBytes int
    n, err := w.Write([]byte("0\r\n"))
    totalBytes += n
    if err != nil { return totalBytes, err }

    if len(trailers) != 0 {
        n, err = w.WriteHeaders(trailers)
        totalBytes += n
        if err != nil { return totalBytes, err }
        return totalBytes, nil
    }

    n, err = w.Write([]byte("\r\n"))
    totalBytes += n
    return totalBytes, err
}
