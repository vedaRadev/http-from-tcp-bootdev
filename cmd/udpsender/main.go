package main

import (
    "net"
    "log"
    "bufio"
    "os"
    "fmt"
)

func main() {
    fmt.Println("attempting to resolve udp addr")

    udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
    if err != nil {
        log.Fatalf("failed to resolve udp addr: %s\n", err.Error())
    }
    fmt.Println("resolved udp addr")

    fmt.Println("attempting to create udp connection")
    udpConnection, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        log.Fatalf("failed to create udp connection: %s\n", err.Error())
    }
    fmt.Println("created udp connection")
    defer udpConnection.Close()

    stdinReader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("> ")
        userInput, err := stdinReader.ReadString('\n')
        if err != nil {
            fmt.Printf("failed to read input: %s\n", err.Error())
            continue
        }
        bytesWritten, err := udpConnection.Write([]byte(userInput))
        if err != nil {
            fmt.Printf("failed to write to udp connection: %s\n", err.Error())
            continue
        }
        fmt.Printf("wrote %d bytes\n", bytesWritten)
    }
}
