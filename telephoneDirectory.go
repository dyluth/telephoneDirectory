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
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//TODO - replace this with a real backend object store
//for the time being a map will be sufficient to show it's working
var datastore list.List

func main() {
	fmt.Println("starting a thing..")
	c := make(chan int)
	go StartServer(c)
	<-c
}

func StartServer(c chan int) {

	datastore := list.New()
	//create some default values in the datastore
	datastore.PushFront(TelephoneEntry{"smith", "bill", "1234567890", "1 road name, town name, city, postcode"})

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

	//set the content type appropriately and stick in the json response
	w.Header().Set("Content-Type", "application/json")

	//pull out the query:
	req.ParseForm()
	query := req.PostFormValue("query")

	//switch on the query to see what we need to provide.. and get it!
	switch query {
	case "list":
		//if "sirname" field persent, return just that ket in an array
		//else return whole directory set
		//add a hardcoded response for now:
		js, _ := json.Marshal([]TelephoneEntry{
			TelephoneEntry{"smith", "bill", "1234567890", "1 road name, town name, city, postcode"},
			TelephoneEntry{"smith", "ben", "1234567890", "1 road name, town name, city, postcode"}})
		w.Write(js)
		fmt.Println("added to header: ", js)
	case "create":
		fmt.Println("number 5")
		//the object should not already exist - if it does return an error
		//now just flow into update case, as the rest is the same
	case "update":
		//replace an existing entry with a new one
	case "remove":
		//the object should already exist - if not return an error
		//remove the object from the datastore

	//case i>7: fmt.Println("is > 7") //will be compile error as type int and bool don't match between case and switch
	default:
		io.WriteString(w, "ooh! Questions!\n")
		//io.WriteString(w, req.Method) //eg "GET"
		//look into the request to see what the query values are

		//req.ParseForm()

		if req.PostFormValue("TEST_ECHO") != "" {
			io.WriteString(w, req.PostFormValue("TEST_ECHO"))
		}
	}

	/*
		js, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		w.Write(js)
	*/
}
