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

func setUpDatabaseConnection(config map[string]string) databasehandling.DBHandler {
	//TODO: maybe just provide connection string from here?
	//Create dbHandler object.
	temp, err := databasehandling.GetDatabaseHandler(config["db_username"], config["db_name"], config["db_password"], config["db_sslmode"])
	if err != nil {
		panic(err)
	}

	//If no errors were thrown, assign this object as global database handler.

	return temp
}

func setUpHTTPHandlers(config map[string]string, db databasehandling.DBHandler) {
	temp, err := strconv.Atoi(config["cookie_life_time"])
	if err != nil {
		log.Panic("Could not parse cookie_life_time value.")
	}

	cookieLifeTime := time.Duration(temp) * time.Second
	handlers.Initialize(config["web_assets_path"], cookieLifeTime, server.MainServerMux, db)
}

func setUpAndRunServer(config map[string]string) {
	val := config["redirect_http_to_https"]

	var redirect bool

	if val == "1" {
		redirect = true
	}

	if redirect {
		serverRedirect := server.GetRedirectServer(config["https_port"], config["http_port"])

		go func() {
			err := serverRedirect.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
	}

	mainServer := server.GetTLSServer(config["https_port"])

	log.Print()
	if config["https_port"] == "443" {
		log.Print("Listen to me at https://" + config["host"])
	} else {
		log.Print("Listen to me at https://" + config["host"] + ":" + config["https_port"])
	}

	err := mainServer.ListenAndServeTLS(config["server_cert"], config["server_key"])

	if err != nil {
		panic(err)
	}

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

	db := setUpDatabaseConnection(config)
	setUpHTTPHandlers(config, db)
	setUpAndRunServer(config)
}
