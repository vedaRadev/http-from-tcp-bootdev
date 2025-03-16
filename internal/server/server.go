package server

import (
    "net"
    "fmt"
    "io"
    "bytes"
    "sync/atomic"
    "strconv"
    "http-from-tcp/internal/headers"
    "http-from-tcp/internal/response"
    "http-from-tcp/internal/request"
)

type HandlerError struct {
    Status response.StatusCode
    Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
    listener net.Listener
    closed atomic.Bool
    handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {
    var server Server
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil { return nil, err }
    server.listener = listener
    server.handler = handler
    go server.listen()
    return &server, nil
}

func (s *Server) Close() error {
    s.closed.Store(true)
    return s.listener.Close()
}

func (s *Server) listen() {
    for {
        connection, err := s.listener.Accept()
        if err != nil {
            if s.closed.Load() { return }
            fmt.Printf("failed to accept connection: %v\n", err.Error())
            continue
        }

        fmt.Println("Accepted connection")
        go s.handle(connection)
    }
}

func (s *Server) handle(conn net.Conn) {
    defer conn.Close()

    parsedRequest, err := request.RequestFromReader(conn)
    responseBody := bytes.NewBuffer([]byte{})
    var handlerError *HandlerError
    if err != nil {
        // TODO handle error
        return
    } else {
        handlerError = s.handler(responseBody, parsedRequest)
    }

    if handlerError == nil {
        // TODO handle error
        _ = response.WriteStatusLine(conn, response.STATUS_OK)
    } else {
        // TODO handle error
        _ = response.WriteStatusLine(conn, handlerError.Status)
        responseBody = bytes.NewBuffer([]byte(handlerError.Message))
    }

    headers := headers.NewHeaders()
    headers["Content-Length"] = strconv.Itoa(responseBody.Len())
    headers["Connection"] = "close"
    headers["Content-Type"] = "text/plain"
    _ = response.WriteHeaders(conn, headers)

    _, _ = conn.Write(responseBody.Bytes())
}
