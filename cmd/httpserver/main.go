package main

import (
    "log"
    "fmt"
    "syscall"
    "strconv"
    "os"
    "os/signal"

    "http-from-tcp/internal/server"
    "http-from-tcp/internal/request"
    "http-from-tcp/internal/response"
    "http-from-tcp/internal/headers"
)

const port = 42069

const bodyTemplate string = "" +
"<html>\n" +
"   <head>\n" +
"       <title>%s</title>\n" +
"   </head>\n" +
"   <body>\n" +
"       <h1>%s</h1>\n" +
"       <p>%s</h1>\n" +
"   </body>\n" +
"</html>\n"

func handle(w *response.Writer, req *request.Request) {
    headers := headers.NewHeaders()
    headers.Add("Connection", "close")
    headers.Add("Content-Type", "text/html")
    var body string
    var statusCode response.StatusCode

    switch(req.RequestLine.RequestTarget) {
        case "/yourproblem": {
            statusCode = response.STATUS_BAD_REQUEST
            body = fmt.Sprintf(
                bodyTemplate,
                "400 Bad Request",
                "Bad Request",
                "Your request honestly kinda sucked.",
            )
        }

        case "/myproblem": {
            statusCode = response.STATUS_INTERNAL_SERVER_ERROR
            body = fmt.Sprintf(
                bodyTemplate,
                "500 Internal Server Error",
                "Internal Server Error",
                "Okay, you know what? This one is on me.",
            )
        }

        default: {
            statusCode = response.STATUS_OK
            body = fmt.Sprintf(
                bodyTemplate,
                "200 OK",
                "Success!",
                "Your request was an absolute banger.",
            )
        }
    }

    headers.Add("Content-Length", strconv.Itoa(len(body)))
    // TODO handle potential write errors
    w.WriteStatusLine(statusCode)
    w.WriteHeaders(headers)
    w.WriteBody([]byte(body))
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
