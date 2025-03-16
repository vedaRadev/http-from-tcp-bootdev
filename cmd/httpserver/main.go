package main

import (
    "log"
    "syscall"
    "io"
    "os"
    "os/signal"

    "http-from-tcp/internal/server"
    "http-from-tcp/internal/request"
)

const port = 42069

func handle(w io.Writer, req *request.Request) *server.HandlerError {
    switch(req.RequestLine.RequestTarget) {
        case "/yourproblem": return &server.HandlerError {
            Status: 400,
            Message: "Your problem is not my problem\n",
        }
        case "/myproblem": return &server.HandlerError {
            Status: 500,
            Message: "Woopsie, my bad\n",
        }
    }

    w.Write([]byte("All good, frfr\n"))
    return nil
}

func main() {
	server, err := server.Serve(port, handle)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
