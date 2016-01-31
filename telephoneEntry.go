package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

/*
this represents an entry in the telephoneDirectory

It contians methods for encoding, decoding, and validating the entry with the following details:
	Surname
	Firstname
	Phone number
	Address (optional)

	JSON Format of entry:
		{surname:"", firstname:"", phone_number:"", address=:""}
*/

type TelephoneEntry struct { //TODO make sure names are JSON compatible
	Surname     string `json:"surname"`
	FirstName   string `json:"firstName"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address,omitempty"` //optional field
}

/*
validates the Entry to ensure that it conforms to specification
*/
func Validate(e TelephoneEntry) bool {
	if e.Surname == "" { //surname cant be blank
		return false
	} else if e.FirstName == "" { //first name cant be blank
		return false
	} else if e.PhoneNumber == "" { //phone number cant be blank
		return false
	} else { //now validate some content

		//check phone number to be numbers and whitespace - cna make this better later
		matched, _ := regexp.MatchString("[\\s\\d]*", e.PhoneNumber)
		if !matched {
			return false
		}
		//TODO validate address field..

	}

	return true
}

func LoadFronJSON(s string) []TelephoneEntry {
	dec := json.NewDecoder(strings.NewReader(s))

	//now create an array that big.
	//TODO - use a variable length slice instead?
	entries := []TelephoneEntry{}

	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for i := 0; dec.More(); i++ {

		// decode an array value (Message)
		var m TelephoneEntry
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		entries = append(entries, m)
		fmt.Printf("%v: %v\n", m.Surname, m.FirstName)
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	return entries
}
