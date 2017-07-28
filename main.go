package main

import "os"

func main() {
	// POC. Refactor please
	os.Chdir(".")
	versionControl().init()
	startQs(findQs("."))
	flushLogs(".")
	versionControl().push()
}
