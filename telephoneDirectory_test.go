package main

import (
	"fmt"
	"net/http"
	"testing"
	"io/ioutil"
)



func TestServer(t *testing.T) {
	c := make(chan int)
	go StartServer(c)
	fmt.Println("testing the server")
	response, err := http.Get("http://localhost:8084/directory")
	fmt.Println("we got: ", response, "and error: ", err)
	if err != nil {
		t.Fatal(err)
	}
	contents, err := ioutil.ReadAll(response.Body)
	
	fmt.Println("contents: ", contents, "and err: ", err)
	if err != nil {
		t.Fatal(err)
	}
	s := string(contents[:])
	fmt.Println("contents: ",  s)
	
	if(s != "ooh! Questions!\n") {
		t.FailNow()
	}	
}

/*
func TestList(t *testing.T) {
	fmt.Println("testing a thing..")
	t.FailNow()
}
*/
