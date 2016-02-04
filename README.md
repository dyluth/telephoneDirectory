# telephoneDirectory
A project written in Go for Cam to learn the language

# brief:
Create a simple HTTP service to represent a phone book.

Acceptance criteria.
- List all entries in the phone book.
- Create a new entry to the phone book.
- Remove an existing entry in the phone book.
- Update an existing entry in the phone book.
- Search for entries in the phone book by surname.

A phone book entry must contain the following details:
- Surname
- Firstname
- Phone number
- Address (optional)

The solution can be in any language. Please upload your project to github and provide us with the URL.
We are not looking for a client or UI for this solution, a simple HTTP based service will suffice.

# Implementation:
telephoneDirectory 		is the main server
telephoneDirectory_test is the test code for the main server - it creates a client and connects to the server to test it
telephoneEntry 			deals with what data is stored in it, validation of the data and converting to and from JSON (and arrags and maps)
telephoneEntry_test		is a simple set of tests for telephoneEntry