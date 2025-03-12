package request

import (
    "io"
    "strings"
    "errors"
    "regexp"
)

type Request struct {
    RequestLine RequestLine
}

type RequestLine struct {
    HttpVersion string
    RequestTarget string
    Method string
}

var versionRegex, _ = regexp.Compile("^\\d\\.\\d$")
var methodRegex, _ = regexp.Compile("^[A-Z]+$")

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
    if len(versionParts) != 2 { return nil, errors.New("invalid protocol/version format") }

    method := requestLineParts[0]
    requestTarget := requestLineParts[1]
    protocol := versionParts[0]
    version := versionParts[1]

    if strings.ToUpper(method) != method { return nil, errors.New("method must be uppercase") }
    if !methodRegex.MatchString(method) { return nil, errors.New("invalid method formaat") }
    if !versionRegex.MatchString(version) { return nil, errors.New("invalid version number format") }
    // NOTE: temporary
    if protocol != "HTTP" { return nil, errors.New("method must be HTTP") }
    // NOTE: temporary
    if version != "1.1" { return nil, errors.New("we only support version 1.1 for now") }

    requestLine.Method = method
    requestLine.RequestTarget = requestTarget
    requestLine.HttpVersion = version

    return &Request { RequestLine: requestLine }, nil
}
