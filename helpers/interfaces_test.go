package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapStructToInterface(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected TestStruct
		hasError bool
	}{
		{
			name: "Valid input",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			expected: TestStruct{
				Name: "John",
				Age:  30,
			},
			hasError: false,
		},
		{
			name: "Invalid input type",
			input: map[string]interface{}{
				"name": "John",
				"age":  "thirty",
			},
			expected: TestStruct{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TestStruct
			err := MapStructToInterface(tt.input, &result)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
