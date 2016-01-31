package main

/*
this
	JSON Format of entry:
		{surname:"", firstname:"", phone_number:"", address=:""}
*/
import (
	"encoding/json"
	"testing"
)

func TestValidate(t *testing.T) {

	validateTests := map[TelephoneEntry]bool{
		TelephoneEntry{"smith", "bill", "1234567890", "1 road name, town name, city, postcode"}:     true,  //good - all optionals included
		TelephoneEntry{"smith", "bill", " 1234 567 890 ", "1 road name, town name, city, postcode"}: true,  //good - valid phone_number spacing
		TelephoneEntry{"smith", "bill", "1234567890", ""}:                                           true,  //good - ommiting address
		TelephoneEntry{"", "bill", "1234567890", ""}:                                                false, //bad firstname - missing
		TelephoneEntry{"smith", "", "1234567890", ""}:                                               false, //bad surname - missing
		TelephoneEntry{"smith", "", "1234as56f7  890", ""}:                                          false, //bad phone number - letters
		TelephoneEntry{"smith", "bill", "", ""}:                                                     false, //bad phone_number - missing
	}

	for key, value := range validateTests {
		result := Validate(key)
		if value != result {
			t.Log("error! validating: ", key, " expected:", value, " got:", result)
			t.Fail()
		}
	}
}

/*
test the conversion to and from JSON
*/
func TestLoadFronJSON(t *testing.T) {
	bilbo := TelephoneEntry{"baggins", "bilbo", "393939", "bag end, Bagshot row, Hobbiton, the Shire"}
	js, _ := json.Marshal([]TelephoneEntry{bilbo})
	jsString := string(js[:])
	t.Log("marshalled entry: ", jsString)
	te := LoadFronJSON(jsString)
	if te[0] != bilbo {
		t.Log("error! validating!")
		t.Log(bilbo)
		t.Log("is not Equal to: ")
		t.Log(te[0])
		t.Fail()
	}
}
