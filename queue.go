package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"os"

	yaml "gopkg.in/yaml.v2"
)

func findQs(dir string) []string {
	// prep
	dir = strings.TrimRight(dir, "/")

	// scan
	files, _ := filepath.Glob(dir + "/*/.q")

	// return
	return files
}

func loadQs(qs []string) {
	var wgQ sync.WaitGroup
	for _, q := range qs {
		wgQ.Add(1)
		go func(q string) {
			defer wgQ.Done()
			newQ(q)
		}(q)
	}
	wgQ.Wait()
}

func newQ(qConfigFile string) {
	// initalize
	q := queue{
		Name:       filepath.Base(filepath.Dir(qConfigFile)),
		TasksDir:   filepath.Dir(qConfigFile),
		ConfigFile: qConfigFile,
		Q:          make(chan task),
		ShutdownQ:  make(chan bool),
	}

	// fetch the file
	contents, configReadErr := ioutil.ReadFile(q.ConfigFile)
	if configReadErr != nil {
		log("fatal", "Unable to read in q config for '"+q.ConfigFile+"'")
		return
	}

	// parse
	unmarshalErr := yaml.Unmarshal([]byte(contents), &q)
	if unmarshalErr != nil {
		log("fatal", "Unable to unmarshal q config for '"+q.ConfigFile+"\n\n"+unmarshalErr.Error())
	}

	// defaults
	if q.WorkerInfo.Count == 0 {
		q.WorkerInfo.Count = 10
	}

	// spin up our workers
	for worker := 0; worker < q.WorkerInfo.Count; worker++ {
		go q.work(worker)
	}

	// create failed/completed folders
	os.Mkdir(q.TasksDir+"/"+q.CompletedQ, 0755)
	os.Mkdir(q.TasksDir+"/"+q.FailedQ, 0755)

	// load up our tasks
	tasks, loadTasksErr := ioutil.ReadDir(q.TasksDir)
	if loadTasksErr != nil {
		log("error", "Error loading tasks for q '"+q.Name+"'")
	}

	for _, taskInfo := range tasks {
		// skip over . files, also directories
		if strings.HasPrefix(taskInfo.Name(), ".") || taskInfo.IsDir() {
			continue
		}

		// Inject our task
		q.Q <- task{
			File: q.TasksDir + "/" + taskInfo.Name(),
			Q:    q.Name,
			CMD:  q.WorkerInfo.CMD,
		}
	}

	// shutdown our workers
	for shutdown := 0; shutdown < q.WorkerInfo.Count; shutdown++ {
		q.ShutdownQ <- true
	}
}

// queue contains information about specific queue
type queue struct {
	// basics
	Name       string `yaml:"name"`
	Desc       string `yaml:"description"`
	ConfigFile string
	TasksDir   string

	// queues to place messages
	Q          chan task
	ShutdownQ  chan bool
	FailedQ    string `yaml:"failed.q"`
	CompletedQ string `yaml:"completed.q"`

	// worker information
	WorkerInfo struct {
		Count int    `yaml:"count"`
		CMD   string `yaml:"cmd"`
	} `yaml:"workers"`
}

func (q *queue) work(id int) {
	for {
		select {
		case <-q.ShutdownQ:
			return
		case task := <-q.Q:
			if task.run() {
				q.complete(task)
			} else {
				q.fail(task)
			}
		}
	}
}

func (q *queue) complete(task task) {
}

func (q *queue) fail(task task) {
}
