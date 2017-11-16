package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	yaml "gopkg.in/yaml.v2"
)

type task struct {
	// name of the task
	Name string
	// file it's located in
	File string
	// Q name it's coming from
	Q string
	// command to run(work)
	CMD string
	// yaml file contents key: value pairings
	Args map[string]string `yaml:",inline"`
	// raw contents of the file
	Contents string
	// verbose mode?
	Verbose bool
	// dry run
	DryRun bool
}

func (t *task) run() bool {
	// init
	t.Name = filepath.Base(t.File)

	// yaml?
	contents, _ := ioutil.ReadFile(t.File)
	t.Contents = string(contents)

	// parse
	yaml.Unmarshal([]byte(contents), &t.Args)

	// at this point we don't care about the errors.
	// the reason, is we are not sure if we need the file contents
	// nor are we sure if we really needed to unmarshal anything
	// the true test will be if any of those args gets used
	// elsewhere
	fns := template.FuncMap{
		"task": taskParams,
	}

	tmpl, parseErr := template.New("params").Funcs(sprig.TxtFuncMap()).Funcs(fns).Parse(t.CMD)
	if parseErr != nil {
		log("failed[template]", t.Q+":"+t.Name, false)
		return false
	}

	// hold our buffer
	cmdParsed := new(bytes.Buffer)
	executionErr := tmpl.Execute(cmdParsed, t)
	if executionErr != nil {
		log("missing/invalid[arguments]", t.Q+":"+t.Name, false)
		return false
	}

	// TODO: HACK, Lets figure out the right way to fix this
	cmdScrubbed := strings.Replace(cmdParsed.String(), "&lt;", "<", -1)

	// setup the task
	cmd := exec.Command("bash", "-c", cmdScrubbed)
	if t.Verbose {
		// send everything out!
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// do not execute the command if it's a dry run
	if !t.DryRun {
		err := cmd.Run()
		if err != nil {
			log("failed", t.Q+":"+t.Name, false)
			return false
		}
	}

	// success! #lifegoals
	log("completed", t.Q+":"+t.Name, false)
	return true
}

func taskParams(m map[string]string, key string) (interface{}, error) {
	val, ok := m[key]
	if !ok {
		return nil, errors.New("missing key " + key)
	}
	return val, nil
}
