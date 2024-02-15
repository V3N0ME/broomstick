package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func main() {

	configFile := flag.String("config", "", "config file")
	shouldSchedule := flag.Bool("schedule", false, "run scheduler ?")
	instance := flag.String("instance", "", "instance to control")
	action := flag.String("action", "ON", "ON/OFF cluster")

	flag.Parse()

	if *configFile == "" {
		panic("config is required")
	}

	filename, _ := filepath.Abs(*configFile)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config Configuration

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	if shouldSchedule != nil && *shouldSchedule {

		NewBroomstick(config).Schedule()

	} else {

		if instance == nil || *instance == "" {
			print("-instance flag is required if schedule is false")
			return
		}

		if action == nil {
			print("-action flag is required if schedule is false")
			return
		}

		if *action != "ON" && *action != "OFF" {
			print("-action flag must either be ON/OFF")
			return
		}

		NewBroomstick(config).Run(RunConfig{
			ID:     *instance,
			Action: *action,
		})
	}
}
