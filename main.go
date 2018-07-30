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
	config["host"] = "localhost"
	config["redirect_http_to_https"] = "0"

	for _, line := range lines {
		fmt.Println(line)
		lineElems := strings.Fields(line)
		config[lineElems[0]] = lineElems[1]
	}

	val := config["redirect_http_to_https"]

	var redirect bool

	if val == "1" {
		redirect = true
	}

	fmt.Println(config)

	core.Initialize(config["db_username"], config["db_name"], config["web_assets_path"])

	core.RunTLS(config["https_port"], config["http_port"], redirect, config["host"], config["server_cert"], config["server_key"])
}
