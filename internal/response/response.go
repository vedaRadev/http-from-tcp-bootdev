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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
    var statusLine string

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
    return err
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
    for key, value := range headers {
        header := fmt.Appendf([]byte{}, "%s: %s\r\n", key, value)
        _, err := w.Write(header)
        if err != nil { return err }
    }

    _, err := w.Write([]byte("\r\n"))
    return err
}
