package main

import (
	golog "log"
	"os"
	"testing"
)

func TestVCSInitCommit(t *testing.T) {
	// cleanup old tests ...
	os.Remove("/tmp/init")
	os.Remove("/tmp/commit")

	vcs := &VCS{
		ConfigFile: "t/q1/.vbQ",
	}

	// test out init
	vcs.init()

	// ok, /tmp/init should exist
	if _, err := os.Stat("/tmp/init"); os.IsNotExist(err) {
		golog.Fatalf("init() did not run properly")
	}

	// run commit
	vcs.commit()

	// test out commit
	if _, err := os.Stat("/tmp/commit"); os.IsNotExist(err) {
		golog.Fatalf("commit() did not run properly")
	}

	// lazy, but make sure the .log exists
	if _, err := os.Stat("t/q1/.log"); os.IsNotExist(err) {
		golog.Fatalf("Expecting .log to exist")
	}
}
