package main

import (
	"fmt"
	"net/http"
	"testing"
	"io/ioutil"
	"net/url"
	"strings"
	
    //"bytes"
)

func TestMain(t *testing.T) {
	c := make(chan int)
	go StartServer(c)

}

func TestServerGet(t *testing.T) {
	t.Log("testing the server")
	
	form := url.Values{}
    form.Set("user", "dave")
	form.Add("cake", "yes please")
	
	//req,err := http.NewRequest("GET", "http://localhost:8084/directory", bytes.NewBufferString(form.Encode()))
	//t.Log("request to send:  ")
	//t.Log(req)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//response, err := http.DefaultClient.Do(req)

	
	response, err := http.DefaultClient.PostForm("http://localhost:8084/directory", form)
	 
	if err != nil {
		t.Fatal(err)
		
	}
    t.Log("status: ",response.Status)
	
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	s := string(contents[:])
	t.Log("contents: ",  s)
	
	//t.Log("time to fail?")
	t.FailNow()
	
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
	fmt.Println("contents: ",  s)
	
	if(!strings.HasPrefix(s, "ooh! Questions!")) {
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
