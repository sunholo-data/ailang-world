package main

import (
	"strings"
	"testing"
)

func TestTaskTitle(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"455", "Task: make `world-publish`'s \"store is locked\" refusal say what to do when the World service holds the store; explain the recovery.", "make `world-publish`'s \"store is locked\" refusal say what to do when…"},
		{"heading", "# Heading", "Heading"},
		{"blanks", "\n \n##  A   heading\ttext\nbody", "A heading text"},
		{"sentence", strings.Repeat("word ", 40) + ". Next.", strings.TrimSpace(strings.Repeat("word ", 14)) + "…"},
		{"word", strings.Repeat("界", 200), strings.Repeat("界", 72) + "…"},
		{"empty", "", ""},
		{"whitespace", " \n\t", ""},
		{"label", "Task:", ""},
		{"sentence-end", "tAsK: Fix this. Then that", "Fix this"},
		{"paren", "Fix this (details)", "Fix this"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := taskTitle([]byte(tc.input)); got != tc.want {
				t.Fatalf("%q, want %q", got, tc.want)
			}
		})
	}
}
