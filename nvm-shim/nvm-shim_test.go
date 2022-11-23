package main

import (
	"fmt"
	"testing"
)

type runs[T, K any] struct {
	input    T
	expected K
}

func TestGlobalInstall(t *testing.T) {
	tests := [...]runs[[]string, bool]{
		{[]string{"install"}, false},
		{[]string{"install", "-g"}, true},
		{[]string{"-g", "install"}, true},
		{[]string{"--help"}, false},
		{[]string{"install", "--location=global"}, true},
		{[]string{"isntall", "--global"}, true},
		{[]string{"isnt", "-g"}, true},
		{[]string{"i", "-g"}, true},
		{[]string{"i", "--global"}, true},
	}

	for _, vals := range tests {
		t.Run(fmt.Sprintf("%s should return %t", vals.input, vals.expected), func(t *testing.T) {
			out := isGlobalInstall(vals.input)

			if out != vals.expected {
				t.Errorf("got %t want %t", out, vals.expected)
			}
		})
	}
}
