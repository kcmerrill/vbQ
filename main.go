package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// Version contains our version information
	Version = "dev"
	// Commit will hold our git commit id
	Commit = "dev"
)

func main() {
	// setup some simple flags
	vbQConfigFileName := flag.String("vbQ", ".vbQ", "The main vbQ configuration filename")
	dir := flag.String("dir", ".", "Directory to run from")
	help := flag.Bool("help", false, "Display vbQ help")
	version := flag.Bool("version", false, "Display vbQ version")
	flag.Parse()

	// do you want help? the version?
	if *help || *version {
		// print out help
		fmt.Println("\nvbQ(" + Version + "#" + Commit + ")")
		if *help {
			fmt.Print("A simple VCS backed queue with built in workers\n\n")
			flag.PrintDefaults()
			fmt.Print("\n")
		}
		fmt.Println()
		os.Exit(0)
	}

	// switch dirs
	os.Chdir(*dir)
	v := NewVCS(*vbQConfigFileName).init()

	// TODO: implement qConfigFileName
	startQs(findQs(".", ".q"))

	// flush the logs
	flushLogs(v.log)

	// commit the goods
	v.commit()
}
