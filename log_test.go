package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFlushLogs(t *testing.T) {
	toWrite := "/tmp/vbqlog.txt"

	os.Remove(toWrite)

	logs = "kc"
	flushLogs(toWrite)

	logs = "merrill"
	flushLogs(toWrite)

	logContents, logError := ioutil.ReadFile(toWrite)

	if logError != nil {
		t.Fatalf("Unable to write to " + toWrite)
	}

	if strings.TrimSpace(string(logContents)) != "kcmerrill" {
		fmt.Print(string(logContents))
		t.Fatalf("Expected 'kcmerrill' to be flushed to the logs")
	}
}

func TestLog(t *testing.T) {
	// reset
	logs = ""
	log("info", "kc was here", false)
	log("bingo", "washisnameo", false)

	if !strings.Contains(logs, "kc was here") {
		t.Fatalf("Expected info to be written")
	}

	if !strings.Contains(logs, "BINGO") {
		t.Fatalf("Expected bingo to be written in all CAPS")
	}
}
