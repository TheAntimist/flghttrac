package main

import (
	"fmt"
	"net/http"
	"time"
)

func PrintRand() {
	fmt.Println("Hello")
}

func createHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
}

func makeAPICall(client *http.Client, req *http.Request, c chan<- *http.Response) {
	res, err := client.Do(req)
	if err == nil {
		c <- res
	} else {
		fmt.Println("Error in fetching data for request: {method:", req.Method, ", URL:", req.URL, "}")
	}
}
