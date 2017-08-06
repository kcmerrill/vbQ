package main

import (
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
	"strings"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	taskFile := "/tmp/testTask"
	resultsFile := "/tmp/testResults"

	// reset the test
	os.Remove(taskFile)
	os.Remove(resultsFile)

	// create a simple taskfile
	ioutil.WriteFile(taskFile, []byte("name: kcmerrill\nemail: kcmerrill@gmail.com"), 0777)

	task := &task{
		File: taskFile,
		CMD:  "echo {{ .Name }} {{ index .Args \"name\" }} {{ index .Args \"email\" }} > " + resultsFile,
	}

	// lets test the goods ...
	if !task.run() {
		golog.Fatalf("The task.run() should not have failed")
	}

	results, resultsErr := ioutil.ReadFile(resultsFile)
	if resultsErr != nil {
		golog.Fatalf("The task did not complete succesfully ...")
	}

	fmt.Println(string(results))
	if strings.TrimSpace(string(results)) != "testTask kcmerrill kcmerrill@gmail.com" {
		golog.Fatalf("The task command does not relfect expected output.")
	}
}
