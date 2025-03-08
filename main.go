package main

import (
    "fmt"
    "os"
    "errors"
    "io"
    "strings"
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
    file, err := os.Open("messages.txt")
    if err != nil {
        fmt.Printf("error opening file: %s\n", err.Error())
        return
    }

    for line := range(getLinesChannel(file)) {
        fmt.Printf("read: %s\n", line)
    }
}
