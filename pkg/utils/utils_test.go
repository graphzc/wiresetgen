package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPointer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input any
	}{
		{
			name:  "String to pointer",
			input: "test",
		},
		{
			name:  "Int to pointer",
			input: 42,
		},
		{
			name:  "Float to pointer",
			input: 3.14,
		},
		{
			name:  "Bool to pointer",
			input: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()

			result := ToPointer(tc.input)

			assert.Equal(tt, tc.input, *result)
		})
	}
}
