package main

import (
	golog "log"
	"testing"
)

func TestFindQs(t *testing.T) {
	qs := findQs("t/", ".vbQ")
	expected := []string{"t/q1/.vbQ", "t/q2/.vbQ"}
	if qs[0] != expected[0] {
		golog.Fatalf("Unable to find the correct queue")
	}

	if qs[1] != expected[1] {
		golog.Fatalf("Unable to find the correct queue(q1)")
	}
}
