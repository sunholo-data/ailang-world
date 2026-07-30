// Minimized reproducer for a Go compiler code-generation regression present in
// go1.26.0 through go1.26.5 (latest stable at 2026-07-30) on darwin/arm64, and
// ABSENT in go1.25.6 / go1.24.9.
//
// Shape: a LOCAL array literal declared inside a function, indexed by a `range`
// index variable, whose element is assigned to a struct string field that is NOT
// at offset 0, in a composite literal that also contains an interface method call
// AFTER the array access, appended to a slice, followed by `break`.
//
// Symptom: the struct's `Field` string is wrong. Depending on struct layout it is
// either empty or carries a CORRUPT string header (nil data pointer with a garbage
// length), in which case merely printing it segfaults.
//
// It is NOT specific to -race: a default `go build` reproduces. `-gcflags=all=-N`
// (optimizations off) makes it correct; `-gcflags=all=-l` (inlining off) does not.
//
// This file exists so the finding survives as an executable artifact rather than
// prose. See ../README.md for provenance and the mission-side consequences.
package main

import "fmt"

type Row struct {
	N      int
	Field  string
	Reason string
}

type Stringer interface {
	Str() string
}

type myStr string

func (s myStr) Str() string { return string(s) }

func worldsShape() []Row {
	var rows []Row
	texts := []string{"w", ""}
	fields := [...]string{"worldRef", "stateRoot"}
	for i, text := range texts {
		if len(text) == 0 {
			v := Stringer(myStr("empty"))
			rows = append(rows, Row{
				Field: fields[i], Reason: v.Str(),
			})
			break
		}
	}
	return rows
}

func main() {
	got := worldsShape()
	if len(got) != 1 {
		fmt.Println("BUG: len(rows) =", len(got))
		return
	}
	r := got[0]
	// Check the length BEFORE printing: on an affected toolchain the string header
	// can be corrupt, and printing it segfaults instead of reporting.
	n := len(r.Field)
	if n > 1000 {
		fmt.Printf("BUG: len(Field)=%d (corrupt string header)\n", n)
		return
	}
	if r.Field != "stateRoot" {
		fmt.Printf("BUG: Field=%q want %q\n", r.Field, "stateRoot")
		return
	}
	fmt.Println("OK")
}
