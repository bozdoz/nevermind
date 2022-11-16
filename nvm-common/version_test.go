package common

import (
	"fmt"
	"testing"
)

type runs[T, K any] struct {
	input    T
	expected K
}

func TestVersion(t *testing.T) {
	tests := [...]runs[string, Version]{
		{"v1.0.0", "1.0.0"},
		{"1.0.0", "1.0.0"},
		{"1.0", "1.0"},
		{"1", "1"},
		{"1.", ""},
		{"V123", "123"},
		{"a123", ""},
		{"w123", ""},
		{"vv123", ""},
		{"", ""},
		{"0", "0"},
		{"0.0.0", "0.0.0"},
		{"0.1.2.3", ""},
		{"v-1.2.3", ""},
		{"-1.2.3", ""},
		{"1.1.9999999", "1.1.9999999"},
		{"one.two.three", ""},
	}

	for _, vals := range tests {
		t.Run(fmt.Sprintf("%q should return %q", vals.input, vals.expected), func(t *testing.T) {
			v, err := GetVersion(vals.input)

			if vals.expected == "" {
				if err == nil {
					t.Errorf("expected error, but got nil")
				} else {
					expectError := versionError(vals.input)
					if err.Error() != expectError.Error() {
						t.Errorf("got %q, want %q", err, expectError)
					}
				}
			} else if v != vals.expected {
				t.Errorf("got %q, want %q", v, vals.expected)
			}
		})
	}
}
