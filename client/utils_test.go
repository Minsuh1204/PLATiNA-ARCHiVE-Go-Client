package client

import "testing"

func TestVersionToString(t *testing.T) {
	version := Version{Major: 1, Minor: 0, Patch: 0}
	if version.String() != "1.0.0" {
		t.Fatal("Version string is not correct")
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		name     string
		v1       Version
		v2       Version
		expected int // 1 for greater, -1 for less, 0 for equal
	}{
		{
			name:     "Equal versions",
			v1:       Version{Major: 1, Minor: 0, Patch: 0},
			v2:       Version{Major: 1, Minor: 0, Patch: 0},
			expected: 0,
		},
		{
			name:     "Major version greater",
			v1:       Version{Major: 2, Minor: 0, Patch: 0},
			v2:       Version{Major: 1, Minor: 0, Patch: 0},
			expected: 1,
		},
		{
			name:     "Major version less",
			v1:       Version{Major: 1, Minor: 0, Patch: 0},
			v2:       Version{Major: 2, Minor: 0, Patch: 0},
			expected: -1,
		},
		{
			name:     "Minor version greater",
			v1:       Version{Major: 1, Minor: 1, Patch: 0},
			v2:       Version{Major: 1, Minor: 0, Patch: 0},
			expected: 1,
		},
		{
			name:     "Minor version less",
			v1:       Version{Major: 1, Minor: 0, Patch: 0},
			v2:       Version{Major: 1, Minor: 1, Patch: 0},
			expected: -1,
		},
		{
			name:     "Patch version greater",
			v1:       Version{Major: 1, Minor: 0, Patch: 1},
			v2:       Version{Major: 1, Minor: 0, Patch: 0},
			expected: 1,
		},
		{
			name:     "Patch version less",
			v1:       Version{Major: 1, Minor: 0, Patch: 0},
			v2:       Version{Major: 1, Minor: 0, Patch: 1},
			expected: -1,
		},
		{
			name:     "Mixed version greater",
			v1:       Version{Major: 1, Minor: 2, Patch: 3},
			v2:       Version{Major: 1, Minor: 2, Patch: 0},
			expected: 1,
		},
		{
			name:     "Minor greater despite patch lower",
			v1:       Version{Major: 1, Minor: 2, Patch: 1},
			v2:       Version{Major: 1, Minor: 1, Patch: 5},
			expected: 1,
		},
		{
			name:     "Minor less despite patch higher",
			v1:       Version{Major: 1, Minor: 1, Patch: 5},
			v2:       Version{Major: 1, Minor: 2, Patch: 1},
			expected: -1,
		},
		{
			name:     "Major greater despite minor/patch lower",
			v1:       Version{Major: 2, Minor: 0, Patch: 0},
			v2:       Version{Major: 1, Minor: 9, Patch: 9},
			expected: 1,
		},
		{
			name:     "Major less despite minor/patch higher",
			v1:       Version{Major: 1, Minor: 9, Patch: 9},
			v2:       Version{Major: 2, Minor: 0, Patch: 0},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Compare(tt.v2)
			if tt.expected == 0 && result != 0 {
				t.Errorf("expected equal, got %d", result)
			} else if tt.expected > 0 && result <= 0 {
				t.Errorf("expected greater, got %d", result)
			} else if tt.expected < 0 && result >= 0 {
				t.Errorf("expected less, got %d", result)
			}
		})
	}
}
