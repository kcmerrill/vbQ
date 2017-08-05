package main

import (
	"io/ioutil"
	"os/exec"

	yaml "gopkg.in/yaml.v2"
)

// NewVCS will return a new VCS instance
func NewVCS(configFile string) *VCS {
	return &VCS{ConfigFile: configFile}
}

// VCS our commands needed for pre/post hooks
type VCS struct {
	ConfigFile string
	Init       string `yaml:"init"`
	Commit     string `yaml:"commit"`
	log        string `yaml:"log"`
}

func (v *VCS) init() *VCS {
	// fetch the file, if not, treat it like a reqular q(no vcs backing)
	contents, _ := ioutil.ReadFile(v.ConfigFile)

	// parse
	unmarshalErr := yaml.Unmarshal([]byte(contents), &v)
	if unmarshalErr != nil {
		log("error", "Unable to parse "+v.ConfigFile+"\n\n"+unmarshalErr.Error(), true)
	}

	// init the repo ... and/or queue
	if v.Init != "" {
		cmd := exec.Command("bash", "-c", v.Init)
		err := cmd.Run()
		if err != nil {
			log("error", "Initilization failed", true)
		}
	}

	return v
}

func (v *VCS) commit() {
	if v.Commit != "" {
		cmd := exec.Command("bash", "-c", v.Commit)
		err := cmd.Run()
		if err != nil {
			log("error", "Commit failed", true)
		}
	}
}
