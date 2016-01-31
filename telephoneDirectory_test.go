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
		{surname = "", firstname="", phone_number="", Address=""}


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

	//"bytes"
)

func TestMain(t *testing.T) {
	c := make(chan int)
	go StartServer(c)

}

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

func TestServer(t *testing.T) {
	//c := make(chan int)
	//go StartServer(c)
	fmt.Println("testing the server")

	response, err := http.Get("http://localhost:8084/directory")
	fmt.Println("we got: ", response, "and error: ", err)
	if err != nil {
		t.Fatal(err)
	}
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}
	s := string(contents[:])
	fmt.Println("contents: ", s)

	if !strings.HasPrefix(s, "ooh! Questions!") {
		fmt.Println("did not get the expected response")
		t.FailNow()
	}
}

/*
func TestList(t *testing.T) {
	fmt.Println("testing a thing..")
	t.FailNow()
}
*/
