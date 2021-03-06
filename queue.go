package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"os"

	yaml "gopkg.in/yaml.v2"
)

func findQs(dir, qFileName string) []string {
	// prep
	dir = strings.TrimRight(dir, "/")

	// scan
	// current directory
	thisLevel, _ := filepath.Glob(dir + "/" + qFileName)
	// first level
	firstLevel, _ := filepath.Glob(dir + "/*/" + qFileName)
	// yeah yeah ... go globs are silly
	secondLevel, _ := filepath.Glob(dir + "/*/*/" + qFileName)

	// combine everything
	files := append(thisLevel, firstLevel...)
	files = append(files, secondLevel...)

	// return
	return files
}

func startQs(qs []string, dryRun bool) bool {
	// setup our wait groups
	wgQ, wgQLock, failures := sync.WaitGroup{}, &sync.Mutex{}, false

	for {
		processed := 0
		// cycle through each configured q
		for _, q := range qs {
			wgQ.Add(1)
			// work concurrently
			go func(q string) {
				defer wgQ.Done()
				// spin up the q
				p, wasFailures := newQ(q, dryRun)
				processed += p
				// if so, we need to return to get correct exit code
				if wasFailures {
					wgQLock.Lock()
					failures = true
					wgQLock.Unlock()
				}
			}(q)
		}
		wgQ.Wait()
		if processed == 0 || failures {
			break
		}
	}

	// return if there were failures
	return failures
}

func newQ(qConfigFile string, dryRun bool) (int, bool) {
	// initalize
	q := queue{
		// the name of the queue(directory)
		Name: filepath.Base(filepath.Dir(qConfigFile)),
		// where the tasks are located
		TasksDir: filepath.Dir(qConfigFile),
		// the name of the config file, usually .q
		ConfigFile: qConfigFile,
		// 'queue' for golang workers
		Q: make(chan task),
		// prep for shutdown
		ShutdownQ: make(chan bool),
		// bringing mutexy back
		lock: &sync.Mutex{},
	}

	// running in the current directory?
	if q.Name == "." {
		// absolute file path
		actualQPath, actualQPathError := filepath.Abs(".")
		if actualQPathError == nil {
			// if we found something, lets use it.
			q.Name = filepath.Base(actualQPath)
		}
	}

	// fetch the file
	contents, configReadErr := ioutil.ReadFile(q.ConfigFile)
	if configReadErr != nil {
		// do not rerun, and exit 1
		log("error", "Unable to read in q config for '"+q.ConfigFile+"'", false)
		return q.tasks, true
	}

	// parse
	unmarshalErr := yaml.Unmarshal([]byte(contents), &q)
	if unmarshalErr != nil {
		log("error", "Unable to unmarshal q config for '"+q.ConfigFile+"\n\n"+unmarshalErr.Error(), false)
		return q.tasks, true
	}

	// if command is empty, bail
	if q.WorkerInfo.CMD == "" {
		log("error", "'command' not set for workers in q '"+q.Name+"'", false)
		return q.tasks, true
	}

	// defaults
	if q.WorkerInfo.Count == 0 {
		q.WorkerInfo.Count = 1
	}

	// where to toss the completed tasks
	if q.QueueInfo.CompletedQ == "" {
		q.QueueInfo.CompletedQ = ".completed"
	}

	// notice how there isn't a finished default Q?
	// that's so the next time vbQ runs, it will rerun the task
	// feel free and set a failed folder in your configuration

	// spin up our workers
	for worker := 0; worker < q.WorkerInfo.Count; worker++ {
		go q.work(worker)
	}

	// create failed/completed folders
	os.Mkdir(q.TasksDir+"/"+q.QueueInfo.CompletedQ, 0755)
	os.Mkdir(q.TasksDir+"/"+q.QueueInfo.FailedQ, 0755)

	// load up our tasks
	tasks, loadTasksErr := ioutil.ReadDir(q.TasksDir)
	if loadTasksErr != nil {
		log("error", "Error loading tasks for q '"+q.Name+"'", true)
		return q.tasks, true
	}

	// go through each of the tasks
	for _, taskInfo := range tasks {
		// skip over . files, also directories and README's
		if strings.ToLower(taskInfo.Name()) == "readme.md" ||
			strings.HasSuffix(taskInfo.Name(), ".template.yml") ||
			strings.HasPrefix(taskInfo.Name(), ".") ||
			taskInfo.Name() == qConfigFile ||
			taskInfo.IsDir() {
			continue
		}

		// Inject our task
		q.Q <- task{
			// the file of task contents
			File: q.TasksDir + "/" + taskInfo.Name(),
			// name of the Q it's coming from
			Q: q.Name,
			// the worker command to run
			CMD: q.WorkerInfo.CMD,
			// any arguments/params the task has given us
			Args: make(map[string]string),
			// verbose mode?
			Verbose: q.WorkerInfo.Verbose,
			// dryrun mode
			DryRun: dryRun,
		}

		// increment the tasks counter
		q.tasks++
	}

	// shutdown our workers
	for shutdown := 0; shutdown < q.WorkerInfo.Count; shutdown++ {
		q.ShutdownQ <- true
	}

	// well?
	return q.tasks, q.wasFailures
}

// queue contains information about specific queue
type queue struct {
	// basics
	// name of the queue
	Name string `yaml:"name"`
	// description of the queue
	Desc        string `yaml:"description"`
	ConfigFile  string
	TasksDir    string
	wasFailures bool
	lock        *sync.Mutex
	tasks       int

	// queues to place messages
	Q         chan task
	ShutdownQ chan bool
	QueueInfo struct {
		FailedQ    string `yaml:"failed"`
		CompletedQ string `yaml:"completed"`
	} `yaml:"queue"`

	// worker information
	WorkerInfo struct {
		// show errors/success?
		Verbose bool   `yaml:"verbose"`
		Count   int    `yaml:"count"`
		CMD     string `yaml:"command"`
	} `yaml:"workers"`
}

func (q *queue) work(id int) {
	// for as long as it takes
	for {
		select {
		// main loop will send a shutdown signal
		case <-q.ShutdownQ:
			return
		// process the task
		case task := <-q.Q:
			// succesfully executed with no errors?
			if task.run() {
				q.complete(task)
			} else {
				// boo!
				q.fail(task)
			}
		}
	}
}

func (q *queue) complete(task task) {
	// completed tasks are simply moved to another directory
	os.Rename(task.File, filepath.Dir(task.File)+"/"+q.QueueInfo.CompletedQ+"/"+task.Name)
}

func (q *queue) fail(task task) {
	// bad status :(
	q.lock.Lock()
	q.wasFailures = true
	q.lock.Unlock()

	// if failedq not setup, tasks will stay to be reprocessed
	if q.QueueInfo.FailedQ != "" {
		os.Rename(task.File, filepath.Dir(task.File)+"/"+q.QueueInfo.FailedQ+"/"+task.Name)
	}
}
