package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var logs string

func flushLogs(filename string) {
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	fmt.Fprintf(file, "%s", logs)
}

func log(t, msg string) {
	logs += fmt.Sprintln(time.Now().Format(time.RFC3339), strings.ToUpper(t), msg)
	fmt.Println(time.Now().Format(time.RFC3339), strings.ToUpper(t), msg)
	if t == "fatal" {
		os.Exit(42)
	}
}
