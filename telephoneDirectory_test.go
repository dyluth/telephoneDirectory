package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	fmt.Println("testing the server")
	response, err := http.Get("http://localhost:8084/directory")
	fmt.Println("we got: ", response)
	if err != nil {
		t.Fatal(err)
	}
	r, req := http.ReadResponse(response, err)
	t.FailNow()
}

func TestList(t *testing.T) {
	fmt.Println("testing a thing..")
	t.FailNow()
}
