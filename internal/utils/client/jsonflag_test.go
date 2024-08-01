package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONFlag_String(t *testing.T) {
	tests := []struct {
		name     string
		target   interface{}
		expected string
	}{
		{
			name:     "string target",
			target:   "hello",
			expected: `"hello"`,
		},
		{
			name:     "integer target",
			target:   123,
			expected: "123",
		},
		{
			name:     "map target",
			target:   map[string]interface{}{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "slice target",
			target:   []int{1, 2, 3},
			expected: `[1,2,3]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &JSONFlag{Target: tt.target}
			result := flag.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJSONFlag_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		err      bool
	}{
		{
			name:     "valid json string",
			input:    `"hello"`,
			expected: "hello",
			err:      false,
		},
		{
			name:     "valid json integer",
			input:    "123",
			expected: float64(123),
			err:      false,
		},
		{
			name:     "valid json map",
			input:    `{"key":"value"}`,
			expected: map[string]interface{}{"key": "value"},
			err:      false,
		},
		{
			name:     "valid json slice",
			input:    `[1,2,3]`,
			expected: []interface{}{float64(1), float64(2), float64(3)},
			err:      false,
		},
		{
			name:     "invalid json",
			input:    `{"key":}`,
			expected: nil,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target interface{}
			flag := &JSONFlag{Target: &target}
			err := flag.Set(tt.input)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if result, ok := target.(map[string]interface{}); ok {
					expected, _ := tt.expected.(map[string]interface{})
					assert.Equal(t, expected, result)
				} else if result, ok := target.([]interface{}); ok {
					expected, _ := tt.expected.([]interface{})
					assert.ElementsMatch(t, expected, result)
				} else {
					assert.Equal(t, tt.expected, target)
				}
			}
		})
	}
}

func TestJSONFlag_Type(t *testing.T) {
	flag := &JSONFlag{}
	result := flag.Type()
	assert.Equal(t, "json", result)
}
