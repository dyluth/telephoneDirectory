package main

/*
Acceptance criteria:

	List all entries in the phone book.
	Create a new entry to the phone book.
	Remove an existing entry in the phone book.
	Update an existing entry in the phone book.
	Search for entries in the phone book by surname.

A phone book entry must contain the following details:
	Surname
	Firstname
	Phone number
	Address (optional)

Thoughts:
	How to uniquely identify a phonebook entry - duplicate names, addresses, - phone number should be unique.. so key off that.
	create a PhoneBookEntry to look like url.Values - easy to encode / decode between places
		need a validate mechanism - to ensure that th phonebookEntry has the right fields, no extras, (NOTE: Address is optional)

	Store this somehow.. possibly just an interface for the moment that records in memory - use a map for the moment
		can replace that with real store later
	need to be able to search on surname only
		can use a dumb brute force search initially

	storing something new (create) and existing (replace) could have the same method, just if it already exists, throw away the previous one
	removing something - should just be a map key removal

	list all entries - need a way to itentify that - some sort of "all" keyword in the POST?

*/

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

func StartServer(c chan int) {

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

	//req.ParseForm()

	if req.PostFormValue("TEST_ECHO") != "" {
		io.WriteString(w, req.PostFormValue("TEST_ECHO"))
	}

	//io.WriteString(w, "form: "+req.Form.Encode())
	fmt.Print("request: ")
	fmt.Println(req)
	fmt.Print("\nform ")
	fmt.Println(req.PostForm)
	fmt.Println(" - ")
	fmt.Println(req.PostFormValue("cake"))

	//in the body, in the vars somewhere..

}
