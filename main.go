package main

import (
    "fmt"
    "os"
    "errors"
    "io"
)

func main() {
    file, err := os.Open("messages.txt")
    if err != nil {
        fmt.Println("Failed to open file")
        return
    }

    bytes := make([]byte, 8, 8)
    for {
        n, err := file.Read(bytes)
        if err != nil {
            if errors.Is(err, io.EOF) { break }
            fmt.Printf("error: %s\n", err.Error())
            break
        }

        fmt.Printf("read: %s\n", bytes[:n])
    }
}
