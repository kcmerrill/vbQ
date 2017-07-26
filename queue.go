package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

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
	var wg sync.WaitGroup
	for _, q := range qs {
		wg.Add(1)
		go func(q string) {
			defer wg.Done()
			newQ(q)
		}(q)
	}
	wg.Wait()
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
			break
		case task := <-q.Q:
			task.run()
		}
	}
}
