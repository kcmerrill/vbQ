package main

import (
	"io/ioutil"
	golog "log"
	"os"
	"testing"

	"github.com/rs/xid"
)

func TestFindQs(t *testing.T) {
	qs := findQs("t/", ".vbQ")
	expected := []string{"t/queues/test_vcs/.vbQ"}
	if qs[0] != expected[0] {
		golog.Fatalf("Unable to find the correct queue")
	}
}

func TestStartQSuccesful(t *testing.T) {
	// cleanup previous runs
	os.RemoveAll("t/queues/basic/.completed")

	// generate a task
	taskID := xid.New().String()
	ioutil.WriteFile("t/queues/basic/"+taskID, []byte(""), 0644)
	ioutil.WriteFile("/tmp/"+taskID, []byte(""), 0644)

	wasFailures := startQs([]string{"t/queues/basic/.q"}, false)
	if wasFailures {
		golog.Fatalf("Failures with running queues.")
	}

	// ok, so two files should have been created
	// .completed/taskname and /tmp/taskname
	// verify both are present
	_, taskWorkerError := os.Stat("/tmp/" + taskID)
	if taskWorkerError != nil {
		golog.Fatalf("[basic] queue worker failed to produce '/tmp/" + taskID)
	}

	_, taskWorkerCompletedError := os.Stat("t/queues/basic/.completed/" + taskID)
	if taskWorkerCompletedError != nil {
		golog.Fatalf("[basic] queue worker failed to complete '.complete/" + taskID)
	}
}

func TestStartQErrorsWithWorkers(t *testing.T) {
	// generate a task
	wasFailures := startQs([]string{"t/queues/errors/.q"}, false)
	if !wasFailures {
		golog.Fatalf("[errors] Expecting errors with startingQs with bad tasks")
	}
}

func TestStartQErrorsWithQConfigMissingCMD(t *testing.T) {
	wasFailures := startQs([]string{"t/queues/errors.with.workers/.q"}, false)
	if !wasFailures {
		golog.Fatalf("[errors.with.workers] Expecting an error when 'command' is not set")
	}
}
