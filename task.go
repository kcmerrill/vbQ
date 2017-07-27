package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type task struct {
	Name string
	File string
	Q    string
	CMD  string
	Args map[string]string `yaml:",inline"`
}

func (t *task) run() bool {
	// init
	t.Name = filepath.Base(t.File)

	// yaml?
	contents, _ := ioutil.ReadFile(t.File)

	// parse
	yaml.Unmarshal([]byte(contents), &t.Args)

	// at this point we don't care about the errors.
	// the reason, is we are not sure if we need the file contents
	// nor are we sure if we really needed to unmarshal anything
	// the true test will be if any of those args gets used
	// elsewhere
	tmpl, parseErr := template.New("params").Parse(t.CMD)
	if parseErr != nil {
		log("failed[template]", t.Q+":"+t.Name)
		return false
	}

	// hold our buffer
	cmdParsed := new(bytes.Buffer)
	executionErr := tmpl.Execute(cmdParsed, t)
	if executionErr != nil {
		log("failed[arguments]", t.Q+":"+t.Name)
		return false
	}

	// actually run the task now
	cmd := exec.Command("bash", "-c", cmdParsed.String())
	err := cmd.Run()
	if err != nil {
		log("failed", t.Q+":"+t.Name)
		return false
	}

	// success! #lifegoals
	log("completed", t.Q+":"+t.Name)
	return true
}
