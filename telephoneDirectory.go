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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

//TODO - replace this with a real backend object store
//for the time being a map will be sufficient to show it's working
var datastore map[string]TelephoneEntry
var datastoreCount int = 137

func main() {
	fmt.Println("starting a thing..")
	c := make(chan int)
	go StartServer(c)
	<-c
}

/*
convenience method to keep the key in sync with the UID field in the telephoneEntry
*/
func addDatastoreEntry(te TelephoneEntry) {
	te.UID = datastoreCount
	datastore[strconv.Itoa(datastoreCount)] = te
	datastoreCount++
}

func listDatastoreEntries() []TelephoneEntry {

	entries := make([]TelephoneEntry, 0, len(datastore))

	for k := range datastore {
		entries = append(entries, datastore[k])
	}
	return entries
}

func StartServer(c chan int) {

	datastore = make(map[string]TelephoneEntry)
	addDatastoreEntry(TelephoneEntry{0, "smith", "bill", "1234567890", "1 road name, town name, city, postcode"})
	addDatastoreEntry(TelephoneEntry{0, "smith", "ben", "987654321", "2 road name, town name, city, postcode"})
	addDatastoreEntry(TelephoneEntry{0, "baggins", "ben", "987654321", "2 road name, town name, city, postcode"})
	//create some default values in the datastore

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
	command := req.PostFormValue("command")

	//fmt.Println("command: ", command)

	//switch on the query to see what we need to provide.. and get it!
	switch command {
	case "list":
		//if "sirname" field persent, return just that ket in an array
		//else return whole directory set
		//add a hardcoded response for now:

		js, err := json.Marshal(datastore)
		if err != nil {
			fmt.Println("ERROR, marsheling: ", err.Error())
			return
		}
		w.Write(js)
		break
	case "create":
		//the object should not already exist - if it does return an error
		//now just flow into update case, as the rest is the same

		break
	case "update":
		//replace an existing entry with a new one
		//now look for "update" - it should contain the JSON to describe the updared person
		entryString := req.PostFormValue("update")
		//check to make sure that it exists
		entry := LoadFromJSON([]byte(entryString))

		//fmt.Println("entry string: ", entryString)

		if _, present := datastore[strconv.Itoa(entry.UID)]; !present {
			//return an error if it does not exist (present will be false)
			//fmt.Println("entry doesnt exist!  ")
			w.WriteHeader(400)
			return
		}
		//once found, simply replace the original with the new one..
		//fmt.Println("loaded ", datastore[strconv.Itoa(entry.UID)])
		datastore[strconv.Itoa(entry.UID)] = entry
		//fmt.Println("updated? ", datastore[strconv.Itoa(entry.UID)])

		break
	case "remove":
		//specify just the ID in a string
		//the object should already exist - if not return an error
		//remove the object from the datastore
		keyToRemove := req.PostFormValue("remove")

		if _, present := datastore[keyToRemove]; !present {
			//return an error if it does not exist (present will be false)
			w.WriteHeader(400)
			return
		}

		delete(datastore, keyToRemove)

		break
	default:
		io.WriteString(w, "ooh! we didnt get any questions!\n")
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
