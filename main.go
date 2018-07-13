package main

import "github.com/d0ku/database_project_go/core"

func main() {
	//TODO: when app is about to be deployed, that config line should be changed.
	core.Initialize("d0ku", "test_database", "./page/")
	core.Run("1234")
}
