package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

var nbOfJumps = 0

func urlTick() (client *http.Client) {
	// create a custom error to know if a redirect happened
	var RedirectAttemptedError = errors.New("redirect")
	client = &http.Client{}
	// return the error, so client won't attempt redirects
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return RedirectAttemptedError
	}
	return client
}

func showResponse(response *http.Response) {
	nbOfJumps++
	if nbOfJumps > 1 {
		fmt.Println("Redirected to ...")
	}
	fmt.Printf("[#%v] %v", nbOfJumps, response.Request.URL.String())
	fmt.Printf("\n Status : %v\n", response.Status)
	for i, v := range response.Header {
		if i != "Location" {
			fmt.Printf(" > %v : %v\n", i, v)
		}
	}
}

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Please, specify an url ...")
		return
	}

	url := args[0]
	client := urlTick()

	for {
		resp, _ := client.Head(url)
		resp.Body.Close()
		showResponse(resp)
		if resp.Header.Get("location") != "" {
			url = resp.Header.Get("location")
		} else {
			break
		}
	}

	nbOfJumps = 0
}
