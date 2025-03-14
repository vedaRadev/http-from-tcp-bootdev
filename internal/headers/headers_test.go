package headers

import (
    "testing"
    // "fmt"
    // "io"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRequestHeaderParse(t *testing.T) {
    // Test: Valid single header
    headers := NewHeaders()
    data := []byte("Host: localhost:42069\r\n\r\n")
    n, done, err := headers.Parse(data)
    require.NoError(t, err)
    require.NotNil(t, headers)
    assert.Equal(t, "localhost:42069", headers.Get("Host"))
    assert.Equal(t, 23, n)
    assert.False(t, done)

    // Test: Invalid spacing header
    headers = NewHeaders()
    data = []byte("       Host : localhost:42069       \r\n\r\n")
    n, done, err = headers.Parse(data)
    require.Error(t, err)
    assert.Equal(t, 0, n)
    assert.False(t, done)

    // Test: Valid done
    headers = NewHeaders()
    data = []byte("\r\n")
    n, done, err = headers.Parse(data)
    require.NotNil(t, headers)
    require.NoError(t, err)
    assert.Equal(t, n, 2)
    assert.True(t, done)

    // Test: 2 valid headers
    headers = NewHeaders()
    n, done, err = headers.Parse([]byte("Accept: application/json\r\n"))
    require.NoError(t, err)
    assert.False(t, done)
    assert.Equal(t, 26, n)
    n, done, err = headers.Parse([]byte("Host: localhost:42069\r\n"))
    require.NoError(t, err)
    assert.False(t, done)
    assert.Equal(t, 23, n)
    assert.Equal(t, "application/json", headers.Get("Accept"))
    assert.Equal(t, "localhost:42069", headers.Get("Host"))

    // Test: missing value
    headers = NewHeaders()
    n, done, err = headers.Parse([]byte("Accept: \r\n"))
    require.Error(t, err)

    // Test: missing name
    headers = NewHeaders()
    n, done, err = headers.Parse([]byte(": value\r\n"))
    require.Error(t, err)

    // Test: illegal field name
    headers = NewHeaders()
    n, done, err = headers.Parse([]byte("@ccept: application/json\r\n"))
    assert.False(t, done)
    assert.Equal(t, n, 0)
    require.Error(t, err)
}

