package request

import (
    "io"
    "strings"
    "errors"
)

type Request struct {
    RequestLine RequestLine
}

type RequestLine struct {
    HttpVersion string
    RequestTarget string
    Method string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
    var requestLine RequestLine

    bytes, err := io.ReadAll(reader)
    if err != nil { return nil, err }

    plaintext := string(bytes)
    parts := strings.Split(plaintext, "\r\n")
    requestLineString := parts[0]
    requestLineParts := strings.Split(requestLineString, " ")
    if len(requestLineParts) != 3 { return nil, errors.New("invalid number of request line parts") }
    versionParts := strings.Split(requestLineParts[2], "/")
    if len(versionParts) != 2 { return nil, errors.New("invalid http version format") }
    version := versionParts[1]

    requestLine.Method = requestLineParts[0]
    requestLine.RequestTarget = requestLineParts[1]
    requestLine.HttpVersion = version

    return &Request { RequestLine: requestLine }, nil
}
