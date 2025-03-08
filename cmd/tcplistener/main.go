package main

import (
    "fmt"
    "log"
    "errors"
    "io"
    "strings"
    "net"
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

func getLinesChannel(f io.ReadCloser) <-chan string {
    linesChan := make(chan string)
    go getLines(f, linesChan)
    return linesChan
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

        for line := range(getLinesChannel(connection)) {
            fmt.Printf("read: %s\n", line)
        }
    }

}
