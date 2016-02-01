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

	JSON Format of entry:
		{surname:"", firstname:"", phone_number:"", address:""}


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
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"encoding/json"
	"strconv"
)

/*
this is needed to start the server in advance of the other tests being run.
the server is started in parallel so that the process doesnt block indefinitely.
*/
func TestMain(t *testing.T) {
	c := make(chan int)
	go StartServer(c)

}

/*
this is just a dumb test to make sure that the telephone directory returns page data correctly
it also injects an element on the page to show that form data is correctly parsed.
*/
func TestServerEchoResponse(t *testing.T) {
	testString := "test string - this should appear on the 2nd line"

	t.Log("testing the server returns a value that we set here")
	form := url.Values{}
	form.Add("TEST_ECHO", testString)

	response, err := http.DefaultClient.PostForm("http://localhost:8084/directory", form)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("status: ", response.Status)

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	s := string(contents[:])
	t.Log("contents: ", s)
	if strings.Contains(s, testString) {
		t.Log("contents as expected")
	} else {
		t.Log("contents unexpected.. looking for: ", testString)
		t.FailNow()
	}
}

/*
even dumber test - this just shows that a page response is created when we query it.
*/
func TestServer(t *testing.T) {
	//c := make(chan int)
	//go StartServer(c)
	fmt.Println("testing the server")

	response, err := http.Get("http://localhost:8084/directory")
	fmt.Println("we got: ", response, "and error: ", err)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}
	s := string(contents[:])

	if !strings.HasPrefix(s, "ooh! Questions!") {
		t.Log("did not get the expected response")
		t.Log("contents: ", s)
		t.FailNow()
	}
}

//	TEST: List all entries in the phone book.
//	TEST: Update an existing entry in the phone book.
//	TEST: Remove an existing entry in the phone book.
func TestListAll(t *testing.T) {

	form := url.Values{}
	//list all the entries
	form.Add("command", "list") //list of "" will return everyone
	form.Add("list", "*") //list of "" will return everyone
	str := sentTestRequest(form, t)
	//get the map of TelephoneEntries
	
	te := LoadMapFromJSON(str) 
	t.Log("list value: ", te)

	//pick one and update it
	//just get an arbitrary entry
	var random TelephoneEntry
	for k := range te {
	    random = te[k]
	    break
	}
	updatedName := "Susan"
	random.FirstName=updatedName
	t.Log("entry set to Susan: ", random.UID)
	form = url.Values{}
	form.Add("query", "update") //list of "" will return everyone
	js, _ := json.Marshal(random)
	form.Add("update", string(js[:])) //send the command to update te[0] to have the name "susan"
	str = sentTestRequest(form, t)
	t.Log("update to susan request: ", string(str[:]))

	//list all agian, confirm that the changes have been made
	form = url.Values{}
	//list all the entries
	form.Add("command", "list") //list of "" will return everyone
	form.Add("list", "*") //list of "" will return everyone
	str = sentTestRequest(form, t)
	//get the map of TelephoneEntries
	te2 := LoadMapFromJSON(str) 
	retrievedName 	 :=te2[strconv.Itoa(random.UID)].FirstName
	
	t.Log("this had bettter be ",updatedName," : ",retrievedName)
		
	if updatedName != retrievedName {
		t.Error("name not changed from ", retrievedName, " to ", updatedName)
		t.FailNow()
	} 
	

	//pick one and remove it
	//confirm that the removed one is not there,
	//and the number returned is 1 less
	//t.FailNow() //TODO finish implementing this test!
}

func sentTestRequest(form url.Values, t *testing.T) []byte{
	response, err := http.DefaultClient.PostForm("http://localhost:8084/directory", form)
	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("contents: ", string(contents[:]))
	return contents
}

//	TEST: Create a new entry to the phone book.
//	TEST: Search for entries in the phone book by surname to find this entry
//	TEST: Search for entries in the phone book by surname to find an empty set of people
func TestSearchCreate(t *testing.T) {
	t.FailNow()
}
