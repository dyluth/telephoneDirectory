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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
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

	if !strings.HasPrefix(s, "ooh! we didnt get any questions!") {
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
	form.Add("command", "list")
	form.Add("list", "*") //list of "*" will return everyone
	str, _ := sentTestRequest(form, t)
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
	random.FirstName = updatedName
	t.Log("entry set to ", updatedName, ": ", random.UID)
	form = url.Values{}
	form.Add("command", "update") //list of "" will return everyone
	js, _ := json.Marshal(random)
	form.Add("update", string(js[:])) //send the command to update te[0] to have the name "susan"
	str, _ = sentTestRequest(form, t)
	t.Log("update to susan request: ", string(str[:]))

	//list all agian, confirm that the changes have been made
	form = url.Values{}
	//list all the entries
	form.Add("command", "list")
	form.Add("list", "*") //list of "*" will return everyone
	str, _ = sentTestRequest(form, t)
	//get the map of TelephoneEntries
	te2 := LoadMapFromJSON(str)
	retrievedName := te2[strconv.Itoa(random.UID)].FirstName

	t.Log("this had bettter be ", updatedName, " : ", retrievedName)

	if updatedName != retrievedName {
		t.Error("name not changed from ", retrievedName, " to ", updatedName)
		t.FailNow()
	}
}

//pick one and remove it
//confirm that the removed one is not there,
//and the number returned is 1 less
//t.FailNow() //TODO finish implementing this test!
func TestRemoveEntry(t *testing.T) {

	form := url.Values{}
	//list all the entries
	form.Add("command", "list")
	form.Add("list", "*") //list of "*" will return everyone
	str, _ := sentTestRequest(form, t)
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

	form = url.Values{}
	//list all the entries
	form.Add("command", "remove")                //tell it that we want to remove something
	form.Add("remove", strconv.Itoa(random.UID)) //fill in the remove field here
	_, code := sentTestRequest(form, t)
	if code != 200 {
		t.Error("ERROR did not get 'OK 200' from request to remove object.  Code: ", code)
		t.FailNow()
	}

	//now list all to make sure that it has been removed
	form = url.Values{}
	//list all the entries
	form.Add("command", "list")
	form.Add("list", "*") //list of "*" will return everyone
	str, _ = sentTestRequest(form, t)
	te = LoadMapFromJSON(str)
	t.Log("list value: ", te)
	for k := range te {
		//check to see if the ID of every Entry returned is the one we have deleted
		if te[k].UID == random.UID {
			//if it is, then the item wasnt really deleted, so fail the test
			t.Error("ERROR: The item ", random.UID, " was NOT deleted as we asked :( ")
			t.FailNow()
		}
	}
}

func TestCreateEntry(t *testing.T) {
	//create a new entry
	form := url.Values{}
	form.Add("command", "create")

	gimli := TelephoneEntry{0, "son of Gloin", "gimli", "10101010101", "3, Dwarrowdelf, Under the Mountain, The Misty Mountains"}
	js, _ := json.Marshal(gimli)
	form.Add("create", string(js[:]))
	_, code := sentTestRequest(form, t)
	if code != 200 {
		t.Error("ERROR got non-200 return code: ", code)
		t.FailNow()
	}

	//check to see that it really is there!
	form = url.Values{}
	form.Add("command", "list")
	form.Add("list", "*") //list of "*" will return everyone
	str, _ := sentTestRequest(form, t)
	//get the map of TelephoneEntries
	te := LoadMapFromJSON(str)

	//now look for gimli in the complete listing of the telephone directory to make sure he is there
	found := false
	for k := range te {
		if gimli.FirstName == te[k].FirstName {
			found = true
			break
		}
	}

	if !found {
		t.Error("ERROR: we added ", gimli.FirstName, ", but couldnt find him in the directory! ")
		t.FailNow()
	}

}

func TestListBySurname(t *testing.T) {
	//list all entres with "surname":"smith"
	surname := "smith"
	form := url.Values{}
	//list all the entries
	form.Add("command", "list") //list of "" will return everyone
	form.Add("list", surname)   //list of "" will return everyone
	str, _ := sentTestRequest(form, t)
	//get the map of TelephoneEntries

	te := LoadMapFromJSON(str)
	//make sure that all the entries returned have that surname
	for k := range te {
		if te[k].Surname != surname {
			t.Error("ERROR: asked for entries of surname: ", surname, " got: ", te[k])
			t.FailNow()
			break
		}
	}

	t.Error("ERROR: not implemented! ")
	t.FailNow()
}

//try to remove an entry that doesnt exist
//this shoudl fail with an errorcode
func TestRemoveMissingEntry(t *testing.T) {
	form := url.Values{}
	form.Add("command", "remove") //update - try to update an existing one.. but where that doesnt exist - should fail
	//this entry doesnt exist -the IDs in the repo start from ID:137
	form.Add("remove", "1")
	_, code := sentTestRequest(form, t)
	if code == 200 {
		t.Error("ERROR got 'OK 200' from request to remove a non existing object ")
		t.FailNow()
	}
}

//test try to update an entry that doesnt exist - should fail
func TestUpdateMissingEntry(t *testing.T) {
	form := url.Values{}
	form.Add("command", "update") //update - try to update an existing one.. but where that doesnt exist - should fail
	js, _ := json.Marshal(TelephoneEntry{2, "baggins", "bilbo", "393939", "bag end, Bagshot row, Hobbiton, the Shire"})
	form.Add("update", string(js[:]))
	_, code := sentTestRequest(form, t)
	if code == 200 {
		t.Error("ERROR got 'OK 200' from request to update non existing object ")
		t.FailNow()
	}
}

func sentTestRequest(form url.Values, t *testing.T) ([]byte, int) {
	response, err := http.DefaultClient.PostForm("http://localhost:8084/directory", form)

	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("contents: ", string(contents[:]), " code: ", response.StatusCode)
	return contents, response.StatusCode
}

//	TEST: Search for entries in the phone book by surname to find an empty set of people
func TestSearchEmpty(t *testing.T) {
	//no-one likes gollum, so he had his name removed from the phonebook
	surname := "gollum"
	form := url.Values{}
	//list all the entries
	form.Add("command", "list")
	form.Add("list", surname) //list of "*" will return everyone
	str, _ := sentTestRequest(form, t)
	//get the map of TelephoneEntries
	te := LoadMapFromJSON(str)

	if len(te) > 0 {
		t.Error("ERROR we looked for an entry that should not be in the directory: ", surname, " but was returned: ", te)
		t.FailNow()
	}
}

func TestDirectoryIDsConsistent(t *testing.T) {
	//list all entres and make sure their keys are the same as UIDs
	form := url.Values{}
	//list all the entries
	form.Add("command", "list")
	form.Add("list", "*") //list of "*" will return everyone
	str, _ := sentTestRequest(form, t)
	//get the map of TelephoneEntries

	te := LoadMapFromJSON(str)
	//make sure that all the entries returned have that surname
	for k := range te {
		if strconv.Itoa(te[k].UID) != k {
			t.Error("ERROR: key ", k, " doesnt match UID in entry: ", te[k])
			t.FailNow()
			break
		}
	}
}
