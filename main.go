package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/d0ku/e_register/core"
)

func main() {
	//read config.cfg, parse it and run server adequately

	configPath := flag.String("configFilePath", "config/config.cfg", "Path to your config file.")

	flag.Parse()

	log.Print("Config file path:" + *configPath)

	content, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	lines = lines[0 : len(lines)-1]

	config := make(map[string]string)
	config["host"] = "localhost"
	config["redirect_http_to_https"] = "0"
	config["cookie_life_time"] = "900"

	for _, line := range lines {
		log.Print(line)
		lineElems := strings.Fields(line)
		config[lineElems[0]] = lineElems[1]
	}

	val := config["redirect_http_to_https"]

	var redirect bool

	if val == "1" {
		redirect = true
	}

	temp, err := strconv.Atoi(config["cookie_life_time"])
	if err != nil {
		log.Panic("Could not parse cookie_life_time value.")
	}

	cookieLifeTime := time.Duration(temp)

	core.Initialize(config["db_username"], config["db_name"], config["web_assets_path"], cookieLifeTime)

	core.RunTLS(config["https_port"], config["http_port"], redirect, config["host"], config["server_cert"], config["server_key"])
}
