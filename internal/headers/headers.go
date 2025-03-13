package headers

import (
    "strings"
    "errors"
)

type Headers map[string]string

func NewHeaders() (Headers) { return make(Headers) }

const CRLF = "\r\n"
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
    stringData := string(data)
    if !strings.Contains(stringData, CRLF) { return 0, false, nil }
    if strings.HasPrefix(stringData, CRLF) { return 2, true, nil }
    header := strings.SplitN(stringData, CRLF, 2)[0]
    headerParts := strings.SplitN(header, ":", 2)
    if len(headerParts) != 2 { return 0, false, errors.New("missing key:value separator ':'") }
    fieldName := strings.TrimSpace(headerParts[0])
    fieldValue := strings.TrimSpace(headerParts[1])
    if fieldName == "" { return 0, false, errors.New("missing field name") }
    if fieldName != headerParts[0] { return 0, false, errors.New("field-name cannot contain whitespace before the separator") }
    if fieldValue == "" { return 0, false, errors.New("missing field value") }
    h[fieldName] = fieldValue
    // include the CRLF we split on in the count
    return len(header) + 2, false, nil
}
