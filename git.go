package main

import "os/exec"

type git struct {
}

func (g git) init() {
	cmd := exec.Command("bash", "-c", "git reset HEAD --hard && git clean -fd")
	cmd.Run()
}

func (g git) push() {
	cmd := exec.Command("bash", "-c", "")
	cmd.Run()
}
