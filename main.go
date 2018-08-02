package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/d0ku/e_register/core/databasehandling"
	"github.com/d0ku/e_register/core/handlers"
	"github.com/d0ku/e_register/core/server"
)

func setUpDatabaseConnection(config map[string]string) {
	//Create dbHandler object.
	temp, err := databasehandling.GetDatabaseHandler(config["db_username"], config["db_name"])
	if err != nil {
		panic(err)
	}

	//If no errors were thrown, assign this object as global database handler.

	databasehandling.DbHandler = temp
}

func setUpHTTPHandlers(config map[string]string) {
	temp, err := strconv.Atoi(config["cookie_life_time"])
	if err != nil {
		log.Panic("Could not parse cookie_life_time value.")
	}

	cookieLifeTime := time.Duration(temp)
	handlers.Initialize(config["web_assets_path"], cookieLifeTime)
}

func setUpServer(config map[string]string) {
	val := config["redirect_http_to_https"]

	var redirect bool

	if val == "1" {
		redirect = true
	}

	server.RunTLS(config["https_port"], config["http_port"], redirect, config["host"], config["server_cert"], config["server_key"])
}

func parseConfigFile(configPath *string, config map[string]string) {
	log.Print("CONFIG_FILE_PATH|" + *configPath)

	content, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	lines = lines[0 : len(lines)-1]

	config["host"] = "localhost"
	config["redirect_http_to_https"] = "0"
	config["cookie_life_time"] = "900"

	for _, line := range lines {
		log.Print("CONFIG|" + line)
		lineElems := strings.Fields(line)
		config[lineElems[0]] = lineElems[1]

	}
}

func main() {
	//Parse config.cfg file and start all services accordingly.

	configPath := flag.String("configFilePath", "config/config.cfg", "Path to your config file.")

	flag.Parse()

	config := make(map[string]string)

	parseConfigFile(configPath, config)

	setUpDatabaseConnection(config)
	setUpHTTPHandlers(config)
	setUpServer(config)
}
