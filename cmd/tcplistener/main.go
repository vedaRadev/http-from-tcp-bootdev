package main

import (
    "fmt"
    "log"
    "errors"
    "io"
    "strings"
    "net"

    requestPkg "http-from-tcp/internal/request"
)

func getLines(f io.ReadCloser, linesChan chan<- string) {
    bytes := make([]byte, 8, 8)
    var current string
    for {
        n, err := f.Read(bytes)
        if err != nil {
            if errors.Is(err, io.EOF) { break }
            fmt.Printf("error: %s\n", err.Error())
            break
        }

        parts := strings.Split(string(bytes[:n]), "\n")
        for i := range(len(parts) - 1) {
            linesChan <- current + parts[i]
            current = ""
        }
        current += parts[len(parts) - 1]
    }

    if len(current) > 0 {
        linesChan <- current
    }

    f.Close()
    close(linesChan)
}

func main() {
    netListener, err := net.Listen("tcp", ":42069")
    if err != nil {
        log.Fatalf("Failed to listen to port 42069: %s\n", err.Error());
    }

    for {
        connection, err := netListener.Accept()
        if err != nil {
            fmt.Printf("Failed to accept connection: %s\n", err.Error())
            continue
        }
        fmt.Println("Accepted connection")

        request, err := requestPkg.RequestFromReader(connection)
        if err != nil {
            fmt.Printf("failed to get request: %v\n", err)
            continue
        }
        fmt.Println("Request line:")
        fmt.Printf("- Method: %v\n", request.RequestLine.Method)
        fmt.Printf("- Target: %v\n", request.RequestLine.RequestTarget)
        fmt.Printf("- Version: %v\n", request.RequestLine.HttpVersion)
    }

}
