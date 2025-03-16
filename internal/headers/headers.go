package headers

import (
    "strings"
    "errors"
    "fmt"
)

type Headers map[string]string

func NewHeaders() (Headers) { return make(Headers) }

func isValidFieldName(fieldName string) bool {
    for _, c := range(fieldName) {
        if  (c < '0' || c > '9') &&
            (c < 'a' || c > 'z') &&
            (c < 'A' || c > 'Z') &&
            c != '!' && c != '#' && c != '$' &&
            c != '%' && c != '&' && c != '\'' &&
            c != '*' && c != '+' && c != '-' &&
            c != '.' && c != '^' && c != '_' &&
            c != '`' && c != '|' && c != '~' {
            return false
        }
    }
    return true
}

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
    if !isValidFieldName(fieldName) { return 0, false, errors.New("field-name contains illegal characters") }
    if fieldValue == "" { return 0, false, errors.New("missing field value") }
    h.Set(fieldName, fieldValue)
    // include the CRLF we split on in the count
    return len(header) + 2, false, nil
}

func (h Headers) Get(fieldName string) string {
    return h[strings.ToLower(fieldName)]
}

func (h Headers) Set(fieldName, fieldValue string) {
    loweredFieldName := strings.ToLower(fieldName)
    existingFieldValue, exists := h[loweredFieldName]
    if !exists {
        h[loweredFieldName] = fieldValue
    } else {
        h[loweredFieldName] = fmt.Sprintf("%s, %s", existingFieldValue, fieldValue)
    }
}
