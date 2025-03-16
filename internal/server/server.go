package server

import (
    "net"
    "fmt"
    "sync/atomic"
)

type Server struct {
    listener net.Listener
    closed atomic.Bool
}

func Serve(port int) (*Server, error) {
    var server Server
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil { return nil, err }
    server.listener = listener
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
            if !s.closed.Load() {
                connection.Close()
                // TODO handle error
                // maybe break out of the loop?
            }
            
            continue
        }

        go s.handle(connection)
    }
}

func (s *Server) handle(conn net.Conn) {
    const response string = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!"
    conn.Write([]byte(response))
    conn.Close()
}
