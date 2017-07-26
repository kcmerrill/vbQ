package main

import "path/filepath"

type task struct {
	Name string
	File string
	CMD  string
}

func (t *task) run() bool {
	// init
	t.Name = filepath.Dir(t.File)

	return true
}
