package core

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func handler(response http.ResponseWriter, request *http.Request) {
	//That implementation is not stable and dangerous, but fun so I will keep it for now.

	//If address is "/" point to index, else try to find adequate file in page folder.

	if request.Method == "GET" {
		url := request.URL
		var path string

		var pageBaseAddress = "./page/"

		if url.String() == "/" {
			path = pageBaseAddress + "index.html"

		} else {
			path = pageBaseAddress + url.Path[1:]
		}

		fileContent, err := ioutil.ReadFile(path)

		if err != nil {
			//Write error html content to response and return it.
			fmt.Println(err)

			page := GetPage()
			page.SetAuthor("d0ku")
			page.SetTitle("File Not Found!")
			page.SetDescription("Displays when file could not be found.")
			page.AddCSSFile("base.css")

			page.SetBody("<h1 class=\"warning\">FILE NOT FOUND!</h1>")

			errorContent := []byte(page.GetHTMLString())

			response.Write(errorContent)
			return
		}
		//Set correct Content-Type in response basing on request.

		response.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		response.Write(fileContent)
	} else {
		//POST
		//call some functinos or do something
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {

	fmt.Println(request.Method)
	if request.Method == "GET" {
		t, _ := template.ParseFiles("page/login.gtpl")
		t.Execute(response, nil)
	} else {
		request.ParseForm()
		//Validate data, and check whether it can be used to log into database.
		fmt.Println(request.Form["username"])
	}
}

func deleteHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(request.URL.String()))
}

//Run starts the server at default port and prints info in console how to connect to it.
func Run() {
	var port = "1234"

	fmt.Println("Listen to me at: https://localhost:" + port)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/", handler)
	err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)

	if err != nil {
		panic(err)
	}
}
