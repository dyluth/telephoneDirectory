package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("starting a thing..")
	c := make(chan int)
	go StartServer(c)
	<-c
}
	
func StartServer(c chan int ) {

	//start the webserver listening on port 8084
	//redirect to DirectoryServer
	//have a channel to keep the main process open until the server fails somehow
	http.HandleFunc("/directory", DirectoryServer)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
	c <- 1
}

func DirectoryServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "ooh! Questions!\n")
	//io.WriteString(w, req.Method) //eg "GET"
	//look into the request to see what the query values are
	req.ParseForm()
	//io.WriteString(w, "form: "+req.Form.Encode())
		fmt.Print("request ")
		fmt.Println(req)
		fmt.Print("form ")
		fmt.Println(req.PostForm)
		fmt.Println(req.PostFormValue("cake"))
	//io.WriteString(w, req)

	//in the body, in the vars somewhere..

}
