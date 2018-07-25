package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/d0ku/database_project_go/core"
)

func main() {
	//read config.cfg, parse it and run server adequately
	//TODO: when app is about to be deployed, that config line should be changed.
	content, err := ioutil.ReadFile("config.cfg")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	lines = lines[0 : len(lines)-1]

	config := make(map[string]string)

	for _, line := range lines {
		fmt.Println(line)
		lineElems := strings.Fields(line)
		config[lineElems[0]] = lineElems[1]
	}

	fmt.Println(config)

	core.Initialize(config["username"], config["database_name"], config["html_templates_path"])
	core.Run(config["port"])
}
