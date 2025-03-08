package main

import (
    "fmt"
    "os"
    "errors"
    "io"
    "strings"
)

func main() {
    file, err := os.Open("messages.txt")
    if err != nil {
        fmt.Printf("error opening file: %s\n", err.Error())
        return
    }

    bytes := make([]byte, 8, 8)
    var current string
    for {
        n, err := file.Read(bytes)
        if err != nil {
            if errors.Is(err, io.EOF) { break }
            fmt.Printf("error: %s\n", err.Error())
            break
        }

        parts := strings.Split(string(bytes[:n]), "\n")
        for i := range(len(parts) - 1) {
            fmt.Printf("read: %s\n", current + parts[i])
            current = ""
        }
        current += parts[len(parts) - 1]
    }

    if len(current) > 0 {
        fmt.Printf("read: %s\n", current)
    }
}
