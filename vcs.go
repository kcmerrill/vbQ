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
	Startup    string `yaml:"startup"`
	Shutdown   string `yaml:"shutdown"`
	Log        string `yaml:"log"`
}

func (v *VCS) startup() *VCS {
	// fetch the file, if not, treat it like a reqular q(no vcs backing)
	contents, _ := ioutil.ReadFile(v.ConfigFile)

	// set default log
	if v.Log == "" {
		v.Log = ".log"
	}

	// parse
	unmarshalErr := yaml.Unmarshal([]byte(contents), &v)
	if unmarshalErr != nil {
		log("error", "Unable to parse "+v.ConfigFile+"\n\n"+unmarshalErr.Error(), true)
	}

	// startup the repo ... and/or queue
	if v.Startup != "" {
		cmd := exec.Command("bash", "-c", v.Startup)
		err := cmd.Run()
		if err != nil {
			log("error", "Initilization failed", true)
		}
	}

	return v
}

func (v *VCS) shutdown() {
	if v.Shutdown != "" {
		cmd := exec.Command("bash", "-c", v.Shutdown)
		err := cmd.Run()
		if err != nil {
			log("error", "Commit failed", true)
		}
	}
}
