package main

import "os"

func main() {
	os.Chdir(".")
	versionControl().init()
	startQs(findQs("."))
	flushLogs(".")
	versionControl().push()
}
