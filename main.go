package main

import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("messages.txt")
    if err != nil {
        fmt.Println("Failed to open file")
        return
    }

    for {
        var bytes [8]byte
        n, err := file.Read(bytes[:])
        if n == 0 {
            break
        }
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            break
        }
        fmt.Printf("read: %s\n", bytes)
    }
}
