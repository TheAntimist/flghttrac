package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)  

func Start(configFile *os.File) {

	var config struct {
		Notifications []string
		Sources       []string
	}

	config_file_bytes, _ := ioutil.ReadAll(configFile)

	err := yaml.Unmarshal(config_file_bytes, &config)

	if err != nil {
		fmt.Println("Unable to parse YAML.")
		os.Exit(1)
	}

	fmt.Println(config)

	//sources := GetSuitableSources(config.Sources)

}

func GetSuitableSources(sources []string) []Source {
	sourceFuncs := make([]Source, len(sources))
	for i := 0; i < len(sources); i++ {
		switch strings.ToLower(sources[i]) {
		case "momondo":
			sourceFuncs = append(sourceFuncs, Momondo{})
			break
		case "skyscanner":
			sourceFuncs = append(sourceFuncs, Skyscanner{})
			break
		case "kiwi":
			sourceFuncs = append(sourceFuncs, Kiwi{})
			break
		}
	}
	return sourceFuncs
}
