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
    "crypto/sha256"

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
    respHeaders := headers.NewHeaders()
    respHeaders.Add("Connection", "close")

    requestTarget := req.RequestLine.RequestTarget
    if strings.HasPrefix(requestTarget, "/httpbin") {
        httpbinPath := strings.TrimPrefix(requestTarget, "/httpbin")
        respHeaders.Add("Content-Type", "text/html")
        respHeaders.Add("Transfer-Encoding", "chunked")
        respHeaders.Add("Trailers", "X-Content-SHA256")
        respHeaders.Add("Trailers", "X-Content-Length")
        w.WriteStatusLine(response.STATUS_OK)
        w.WriteHeaders(respHeaders)

        resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", httpbinPath))
        if err != nil {
            // TODO send 400/500 response
            fmt.Printf("Failed to get httpbin stuff: %s\n", err.Error())
            return
        }

        buf := make([]byte, 1024, 1024)
        body := []byte{}
        for {
            n, err := resp.Body.Read(buf)
            if err != nil {
                if !errors.Is(err, io.EOF) {
                    fmt.Printf("error reading httpbin response: %s\n", err.Error())
                }
                break;
            }
            fmt.Printf("read %v bytes\n", n)
            body = append(body, buf[:n]...)
            w.WriteChunkedBody(buf[:n])
        }

        sha := sha256.Sum256(body)
        trailers := headers.NewHeaders()
        trailers.Add("X-Content-SHA256", fmt.Sprintf("%x", sha[:]))
        trailers.Add("X-Content-Length", strconv.Itoa(len(body)))

        w.WriteChunkedBodyDone(trailers)
    } else {
        var body string
        var statusCode response.StatusCode
        switch(requestTarget) {
            case "/video": {
                respHeaders.Add("Content-Type", "video/mp4")
                statusCode = response.STATUS_OK
                data, err := os.ReadFile("./assets/vim.mp4")
                if err != nil {
                    fmt.Printf("failed to read file: %s\n", err.Error())
                }
                body = string(data)
            }

            case "/yourproblem": {
                respHeaders.Add("Content-Type", "text/html")
                statusCode = response.STATUS_BAD_REQUEST
                body = fmt.Sprintf(
                    bodyTemplate,
                    "400 Bad Request",
                    "Bad Request",
                    "Your request honestly kinda sucked.",
                )
            }

            case "/myproblem": {
                respHeaders.Add("Content-Type", "text/html")
                statusCode = response.STATUS_INTERNAL_SERVER_ERROR
                body = fmt.Sprintf(
                    bodyTemplate,
                    "500 Internal Server Error",
                    "Internal Server Error",
                    "Okay, you know what? This one is on me.",
                )
            }

            default: {
                respHeaders.Add("Content-Type", "text/html")
                statusCode = response.STATUS_OK
                body = fmt.Sprintf(
                    bodyTemplate,
                    "200 OK",
                    "Success!",
                    "Your request was an absolute banger.",
                )
            }
        }

        respHeaders.Add("Content-Length", strconv.Itoa(len(body)))
        w.WriteStatusLine(statusCode)
        w.WriteHeaders(respHeaders)
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
