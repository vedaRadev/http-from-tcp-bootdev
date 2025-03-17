package main

import (
    "log"
    "io"
    "fmt"
    "errors"
    "syscall"
    "strconv"
    "strings"
    "net/http"
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

    requestTarget := req.RequestLine.RequestTarget
    if strings.HasPrefix(requestTarget, "/httpbin") {
        httpbinPath := strings.TrimPrefix(requestTarget, "/httpbin")
        headers.Add("Transfer-Encoding", "chunked")
        w.WriteStatusLine(response.STATUS_OK)
        w.WriteHeaders(headers)

        resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", httpbinPath))
        if err != nil {
            // TODO send 400/500 response
            fmt.Printf("Failed to get httpbin stuff: %s\n", err.Error())
            return
        }

        buf := make([]byte, 1024, 1024)
        for {
            n, err := resp.Body.Read(buf)
            if err != nil {
                if !errors.Is(err, io.EOF) {
                    fmt.Printf("error reading httpbin response: %s\n", err.Error())
                }
                break;
            }
            w.WriteChunkedBody(buf[:n])
        }
        w.WriteChunkedBodyDone()
    } else {
        var body string
        var statusCode response.StatusCode
        switch(requestTarget) {
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
        w.WriteStatusLine(statusCode)
        w.WriteHeaders(headers)
        w.WriteBody([]byte(body))
    }
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
