package provider

import (
	"errors"
	"testing"
)

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "not found lowercase",
			err:      errors.New("[ERROR]: not found"),
			expected: true,
		},
		{
			name:     "not_found with underscore",
			err:      errors.New("[ERROR]: not_found"),
			expected: true,
		},
		{
			name:     "Not Found mixed case",
			err:      errors.New("[ERROR]: Not Found"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("[ERROR]: rate limit exceeded"),
			expected: false,
		},
		{
			name:     "empty error",
			err:      errors.New(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNotFoundError(tt.err)
			if result != tt.expected {
				t.Errorf("isNotFoundError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}
