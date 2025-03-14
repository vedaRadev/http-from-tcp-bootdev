package request

import (
    "io"
    "strings"
    "errors"
    "regexp"
    "http-from-tcp/internal/headers"
)

const bufferSize int = 8

const (
    REQ_PARSER_REQUESTLINE = int(iota)
    REQ_PARSER_HEADERS
    REQ_PARSER_DONE
)

type Request struct {
    RequestLine RequestLine
    Headers headers.Headers
    parserState int
}

type RequestLine struct {
    HttpVersion string
    RequestTarget string
    Method string
}

func (requestLine *RequestLine) parse(bytes []byte) (int, error) {
    var requestLineBytes int
    plaintext := string(bytes)
    if !strings.Contains(plaintext, "\r\n") { return 0, nil }

    var method, requestTarget, protocol, version string
    parts := strings.Split(plaintext, "\r\n")
    requestLineString := parts[0]
    requestLineBytes = len(requestLineString) + 2 // +2 for \r\n
    requestLineParts := strings.Split(requestLineString, " ")
    if len(requestLineParts) != 3 { return 0, errors.New("invalid number of request line parts") }
    versionParts := strings.Split(requestLineParts[2], "/")
    if len(versionParts) != 2 { return 0, errors.New("invalid protocol/version format") }

    method = requestLineParts[0]
    requestTarget = requestLineParts[1]
    protocol = versionParts[0]
    version = versionParts[1]

    if strings.ToUpper(method) != method { return 0, errors.New("method must be uppercase") }
    if !methodRegex.MatchString(method) { return 0, errors.New("invalid method formaat") }
    if !versionRegex.MatchString(version) { return 0, errors.New("invalid version number format") }
    // NOTE: temporary
    if protocol != "HTTP" { return 0, errors.New("method must be HTTP") }
    // NOTE: temporary
    if version != "1.1" { return 0, errors.New("we only support version 1.1 for now") }

    requestLine.Method = method
    requestLine.RequestTarget = requestTarget
    requestLine.HttpVersion = version
    return requestLineBytes, nil
}

var versionRegex, _ = regexp.Compile("^\\d\\.\\d$")
var methodRegex, _ = regexp.Compile("^[A-Z]+$")

func (r *Request) parse(data []byte) (int, error) {
    switch (r.parserState) {

    case REQ_PARSER_DONE: 
        return 0, errors.New("tried to read data in done state")

    case REQ_PARSER_REQUESTLINE:
        bytesParsed, err := r.RequestLine.parse(data)
        if err != nil { return 0, err }
        if bytesParsed == 0 { return 0, nil }
        r.parserState = REQ_PARSER_HEADERS
        return bytesParsed, nil

    case REQ_PARSER_HEADERS:
        bytesParsed, done, err := r.Headers.Parse(data)
        if err != nil { return 0, err }
        if bytesParsed == 0 { return 0, nil }
        if done { r.parserState = REQ_PARSER_DONE }
        return bytesParsed, nil

    default:
        return 0, errors.New("unreconized parser state")
    }
}

func RequestFromReader(reader io.Reader) (*Request, error) {
    buf := make([]byte, bufferSize, bufferSize)
    readToIndex := 0
    var request Request
    request.parserState = REQ_PARSER_REQUESTLINE
    request.Headers = headers.NewHeaders()

    for request.parserState != REQ_PARSER_DONE {
        // Grow if full
        if readToIndex == len(buf) {
            newLen := len(buf) * 2
            grown := make([]byte, newLen, newLen)
            copy(grown, buf)
            buf = grown
        }

        bytesRead, err := reader.Read(buf[readToIndex:])
        if errors.Is(io.EOF, err) {
            request.parserState = REQ_PARSER_DONE
            break
        }
        readToIndex += bytesRead

        bytesParsed, err := request.parse(buf[:readToIndex])
        if err != nil { return nil, err }
        if bytesParsed > 0 {
            // remove parsed bytes and shrink
            if bytesParsed == readToIndex {
                buf = make([]byte, bufferSize, bufferSize)
                readToIndex = 0
            } else {
                leftover := readToIndex - bytesParsed
                newLen := max(readToIndex - bytesParsed, bufferSize)
                shrunk := make([]byte, newLen, newLen)
                copy(shrunk, buf[bytesParsed:])
                readToIndex = leftover
                buf = shrunk
            }
        }
    }

    return &request, nil
}
