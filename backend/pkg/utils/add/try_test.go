package add_test

import (
	"testing"

	"spyal/pkg/utils/add"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"2 + 3", 2, 3, 5},
		{"0 + 0", 0, 0, 0},
		{"-1 + 1", -1, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := add.Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
