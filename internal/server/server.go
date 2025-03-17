package server

import (
    "net"
    "fmt"
    "sync/atomic"
    "http-from-tcp/internal/response"
    "http-from-tcp/internal/request"
)

type Handler func(w *response.Writer, req *request.Request)

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
    // TODO handle error
    if err != nil { return }
    responseWriter := response.NewWriter(conn)
    s.handler(&responseWriter, parsedRequest)
}
