package main

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
	Address     string `json:"address,omitempty"`  //optional field
}

/*
validates the Entry to ensure that it conforms to specification
*/
func Validate(e TelephoneEntry) bool {
	return false
}
