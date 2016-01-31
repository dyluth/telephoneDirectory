package main

/*
this
	JSON Format of entry:
		{surname:"", firstname:"", phone_number:"", address=:""}
*/
import (
	"testing"
)

func TestValidate(t *testing.T) {

	validateTests := map[TelephoneEntry]bool{
		TelephoneEntry{"smith", "bill", "1234567890", "1 road name, town name, city, postcode"}: true,
		TelephoneEntry{"smith", "bill", "1234567890",""}:                                           true,
		TelephoneEntry{"", "bill", "1234567890",""}:                                                false,
		TelephoneEntry{"smith", "", "1234567890",""}:                                               false,
		TelephoneEntry{"smith", "bill", "",""}:                                                     false,
	}
	
	for key, value := range validateTests {
		result :=Validate(key)
		if value != result {
			t.Log("error! validating: ",key, " expected:", value, " got:", result)
			t.Fail()
		}
	}

}
