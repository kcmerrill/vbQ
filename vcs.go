package main

func versionControl() vc {
	// use git, hardcoded ... for now
	return git{}
}

type vc interface {
	init()
	push()
}
