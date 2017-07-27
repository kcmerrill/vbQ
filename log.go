package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func log(t, msg string) {
	fmt.Println((time.Now().Format(time.RFC3339)), strings.ToUpper(t), msg)
	if t == "fatal" {
		os.Exit(42)
	}
}
