package main

import "github.com/d0ku/database_project_go/core"

func main() {
	core.Initialize("d0ku", "database_project_go", "./page/")
	core.Run("1234")
}
