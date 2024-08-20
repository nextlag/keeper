package errs

import (
	"bytes"
	"compress/gzip"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePostgresErr(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected GormErr
	}{
		{
			name:     "valid error",
			input:    errors.New(`{"Code":"123","Message":"Test error"}`),
			expected: GormErr{},
		},
		{
			name:     "invalid JSON",
			input:    errors.New("invalid error format"),
			expected: GormErr{},
		},
		{
			name:     "empty error",
			input:    errors.New(""),
			expected: GormErr{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParsePostgresErr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseServerError(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected string
	}{
		{
			name:     "valid JSON error",
			body:     []byte(`{"error":"Test server error"}`),
			expected: "Test server error",
		},
		{
			name:     "valid gzipped JSON error",
			body:     gzipCompress([]byte(`{"error":"Gzipped server error"}`)),
			expected: "Gzipped server error",
		},
		{
			name:     "invalid JSON",
			body:     []byte(`invalid json`),
			expected: "invalid json",
		},
		{
			name:     "empty JSON",
			body:     []byte(`{}`),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseServerError(tt.body)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsGzipped(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected bool
	}{
		{
			name:     "gzipped body",
			body:     gzipCompress([]byte("test")),
			expected: true,
		},
		{
			name:     "non-gzipped body",
			body:     []byte("test"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isGzipped(tt.body)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnzip(t *testing.T) {
	tests := []struct {
		name      string
		body      []byte
		expected  string
		wantError bool
	}{
		{
			name:      "valid gzipped body",
			body:      gzipCompress([]byte("test")),
			expected:  "test",
			wantError: false,
		},
		{
			name:      "invalid gzipped body",
			body:      []byte("invalid"),
			expected:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unzip(tt.body)
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.expected, string(result))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, string(result))
			}
		})
	}
}

// gzipCompress is a helper function to compress data for testing.
func gzipCompress(data []byte) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil
	}
	writer.Close()
	return buf.Bytes()
}
