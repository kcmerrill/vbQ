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
		ConfigFile: "t/queues/test_vcs/.vbQ",
	}

	// test out startup
	vcs.startup()

	// ok, /tmp/vcs_test_startup should exist
	if _, err := os.Stat("/tmp/test_vcs_startup"); os.IsNotExist(err) {
		golog.Fatalf("startup() did not run properly")
	}

	// run shutdown
	vcs.shutdown()

	// test out shutdown
	if _, err := os.Stat("/tmp/test_vcs_shutdown"); os.IsNotExist(err) {
		golog.Fatalf("shutdown() did not run properly")
	}

	// Make sure the .log exists
	if _, err := os.Stat("t/queues/test_vcs/.log"); os.IsNotExist(err) {
		golog.Fatalf("Expecting .log to exist")
	}
}
