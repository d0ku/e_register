package core

import (
	"fmt"
	"io"
	"net/http"
)

func handler(response http.ResponseWriter, request *http.Request) {
	//	fmt.Println(request.TLS)
	io.WriteString(response, "Hello World!")
}

func Run() {
	var port = "1234"

	fmt.Println("Listen to me at: https://localhost:" + port)

	http.HandleFunc("/", handler)
	err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)

	if err != nil {
		panic(err)
	}
}
