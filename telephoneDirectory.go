package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("starting a thing..")
	StartServer()
}

func StartServer() {
	//start the webserver listening on port 8084
	//redirect to DirectoryServer
	http.HandleFunc("/directory", DirectoryServer)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func DirectoryServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "ooh! Questions!\n")
}
